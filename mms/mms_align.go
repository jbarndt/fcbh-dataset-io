package mms

import (
	"bufio"
	"context"
	"dataset"
	"dataset/db"
	"dataset/generic"
	"dataset/input"
	log "dataset/logger"
	"dataset/timestamp"
	"encoding/json"
	"github.com/divan/num2words"
	"golang.org/x/text/unicode/norm"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type MMSAlign_Input struct {
	AudioFile string   `json:"audio_file"`
	NormWords []string `json:"words"`
}

type Word struct {
	scriptId int64
	wordId   int64
	wordSeq  int
	word     string
	uroman   string
}

type MMSAlign struct {
	ctx     context.Context
	conn    db.DBAdapter // This database adapter must contain the text to be processed
	lang    string
	sttLang string
}

func NewMMSAlign(ctx context.Context, conn db.DBAdapter, lang string, sttLang string) MMSAlign {
	var m MMSAlign
	m.ctx = ctx
	m.conn = conn
	m.lang = lang
	m.sttLang = sttLang
	return m
}

// ProcessFiles will perform Forced Alignment on these files
func (a *MMSAlign) ProcessFiles(files []input.InputFile) dataset.Status {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_align.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_FA_PYTHON`), pythonScript, lang)
	if status.IsErr {
		return status
	}
	for _, file := range files {
		log.Info(a.ctx, "MMS Align", file.BookId, file.Chapter)
		status = a.processFile(file, writer, reader)
		if status.IsErr {
			return status
		}
	}
	return status
}

// processFile will process one audio file through mms forced alignment
func (m *MMSAlign) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) dataset.Status {
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_fa_")
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var faInput MMSAlign_Input
	faInput.AudioFile, status = timestamp.ConvertMp3ToWav(m.ctx, tempDir, file.FilePath())
	if status.IsErr {
		return status
	}
	var wordList []Word
	faInput.NormWords, wordList, status = m.prepareText(m.lang, file.BookId, file.Chapter)
	if status.IsErr {
		return status
	}
	content, err := json.Marshal(faInput)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error marshalling json`)
	}
	// development
	//err2 := os.WriteFile("engweb_fa_inp.json", content, 0644)
	//if err2 != nil {
	//	panic(err2)
	//}
	_, err = writer.WriteString(string(content) + "\n")
	if err != nil {
		return log.Error(m.ctx, 500, err, "Error writing to mms_align.py")
	}
	err = writer.Flush()
	if err != nil {
		return log.Error(m.ctx, 500, err, "Error flush to mms_align.py")
	}
	response, err2 := reader.ReadString('\n')
	if err2 != nil {
		return log.Error(m.ctx, 500, err2, `Error reading mms_align.py response`)
	}
	m.processPyOutput(file, wordList, response)
	// development
	//fmt.Println(len(wordList)) // temp
	//err = os.WriteFile("engweb_fa_out.json", []byte(response), 0644)
	//if err != nil {
	//	panic(err)
	//}
	return status
}

func (m *MMSAlign) prepareText(lang string, bookId string, chapter int) ([]string, []Word, dataset.Status) {
	var textList []string
	var refList []Word
	var status dataset.Status
	var dbWords []db.Word
	dbWords, status = m.conn.SelectWordsByBookChapter(bookId, chapter)
	if status.IsErr {
		return textList, refList, status
	}
	for _, word := range dbWords {
		cleanWd := norm.NFC.String(word.Word)
		cleanWd = m.cleanText(cleanWd)
		results := strings.FieldsFunc(cleanWd, func(r rune) bool { // split on hyphen
			return r == '\u002D' || (r >= '\u2010' && r <= '\u2014')
		})
		// Because the parts are all given the same wordId, they are treated as the same
		// word.  Simply discarding th hyphen would do the same thing.
		for _, part := range results {
			var ref Word
			ref.scriptId = int64(word.ScriptId)
			ref.wordId = int64(word.WordId)
			ref.wordSeq = word.WordSeq
			ref.word = part
			refList = append(refList, ref)
			textList = append(textList, strings.ReplaceAll(part, "\u2019", "'"))
		}
	}
	uRoman, status2 := URoman(m.ctx, lang, textList)
	for i := range uRoman {
		uRoman[i] = strings.ToLower(uRoman[i])
	}
	if status2.IsErr {
		return textList, refList, status2
	}
	if len(uRoman) != len(refList) {
		status = log.ErrorNoErr(m.ctx, 500, "Word count did not match in MMS_FA prepareText", bookId, chapter, refList[0].scriptId)
		return textList, refList, status
	}
	textList = nil
	for i := range refList {
		textList = append(textList, uRoman[i])
		refList[i].uroman = uRoman[i]
		if utf8.RuneCountInString(refList[i].word) > utf8.RuneCountInString(uRoman[i]) {
			status = log.ErrorNoErr(m.ctx, 500, "Character count did not match in MMS_FA prepareText", bookId, chapter, refList[i].word, uRoman[i])
		}
	}
	return textList, refList, status
}

func (m *MMSAlign) cleanText(text string) string {
	var result []rune
	for _, ch := range []rune(text) {
		if unicode.IsLetter(ch) || unicode.IsSpace(ch) {
			result = append(result, ch)
		} else if unicode.IsDigit(ch) {
			num := []rune(num2words.Convert(int(ch) - 48))
			result = append(result, num...)
		} else if ch == '\u0027' || ch == '\u2019' {
			result = append(result, '\u0027') // replace any Apostrophe with std one
		} else if ch == '\u002D' || (ch >= '\u2010' && ch <= '\u2014') { // hyphen
			result = append(result, ch)
		} else {
			log.Warn(m.ctx, "Discarded Char in mms_fa.cleanText", string(ch), ch)
		}
	}
	return string(result)
}

type MMSAlignResult struct {
	Ratio      float64        `json:"ratio"`
	Dictionary map[string]int `json:"dictionary"`
	Tokens     [][][]float64  `json:"tokens"`
}

func (m *MMSAlign) processPyOutput(file input.InputFile, wordRefs []Word, response string) dataset.Status {
	var status dataset.Status
	response = strings.TrimRight(response, "\n")
	var mmsAlign MMSAlignResult
	err := json.Unmarshal([]byte(response), &mmsAlign)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error unmarshalling json`)
	}
	var tokenDict = make(map[int]rune)
	for chr, token := range mmsAlign.Dictionary {
		tokenDict[token] = []rune(chr)[0]
	}
	var ok bool
	var faWords [][]generic.Char
	for _, wd := range mmsAlign.Tokens {
		var word []generic.Char
		for _, ch := range wd {
			var char generic.Char
			char.Token = int(ch[0])
			char.Start = ch[1] * mmsAlign.Ratio
			char.End = ch[2] * mmsAlign.Ratio
			char.Score = ch[3]
			char.Uroman, ok = tokenDict[char.Token]
			if !ok {
				log.Warn(m.ctx, "Character not found in tokenDict", char.Token)
			}
			word = append(word, char)
		}
		faWords = append(faWords, word)
	}
	//fmt.Println(response)
	if len(faWords) != len(wordRefs) {
		return log.ErrorNoErr(m.ctx, 400, "Num words input to mms_align:", len(wordRefs), ", num timestamps returned:", len(faWords))
	}
	var words []db.Audio
	for i, ref := range wordRefs {
		var word db.Audio
		word.BookId = file.BookId
		word.ChapterNum = file.Chapter
		word.AudioFile = file.Filename
		word.ScriptId = ref.scriptId
		word.WordId = ref.wordId // because hypenated words were split, multiple words can have the same wordId
		word.WordSeq = ref.wordSeq
		word.Text = ref.word
		word.Uroman = ref.uroman
		faWd := faWords[i]
		word.BeginTS = faWd[0].Start
		word.EndTS = faWd[len(faWd)-1].End
		wordChars := []rune(ref.word)
		uromanChars := []rune(ref.uroman)
		for j, ch := range faWd {
			word.FAScore += ch.Score
			faWd[j].Seq = j
			if j < len(wordChars) { // This assumes uroman can be longer, but not shorter than source word
				faWd[j].Norm = wordChars[j]
			}
			if faWd[j].Uroman != uromanChars[j] {
				log.ErrorNoErr(m.ctx, 500, "Norm", ref.uroman, "does not match")
			}
		}
		word.FAScore = word.FAScore / float64(len(faWd))
		word.Chars = faWd
		words = append(words, word)
	}
	var wordsByLine [][]db.Audio
	wordsByLine = m.groupByLine(words)
	var verses []db.Audio
	verses = m.summarizeByVerse(wordsByLine)
	verses = m.midPoint(verses)
	status = m.conn.UpdateScriptFATimestamps(verses)
	if status.IsErr {
		return status
	}
	status = m.conn.UpdateWordFATimestamps(words)
	if status.IsErr {
		return status
	}
	status = m.conn.InsertAudioChars(words)
	return status
}

func (a *MMSAlign) groupByLine(words []db.Audio) [][]db.Audio {
	var result [][]db.Audio
	if len(words) == 0 {
		return result
	}
	currWd := words[0].ScriptId
	start := 0
	for i, wd := range words {
		if wd.ScriptId != currWd {
			currWd = wd.ScriptId
			oneLine := make([]db.Audio, i-start)
			copy(oneLine, words[start:i])
			result = append(result, oneLine)
			start = i
		}
	}
	if start != len(words) {
		lastLine := make([]db.Audio, len(words)-start)
		copy(lastLine, words[start:])
		result = append(result, lastLine)
	}
	return result
}

func (m *MMSAlign) summarizeByVerse(chapter [][]db.Audio) []db.Audio {
	var result []db.Audio
	for _, verse := range chapter {
		var vs = verse[0]
		vs.EndTS = verse[len(verse)-1].EndTS
		var text []string
		var uroman []string
		var scores []float64
		for _, word := range verse {
			text = append(text, word.Text)
			uroman = append(uroman, word.Uroman)
			scores = append(scores, word.FAScore)
		}
		vs.Text = strings.Join(text, " ")
		vs.Uroman = strings.Join(uroman, " ")
		vs.FAScore = m.average(scores, 5)
		result = append(result, vs)
	}
	return result
}

func (m *MMSAlign) average(scores []float64, precision int) float64 {
	var sum float64
	for _, scr := range scores {
		sum += scr
	}
	avg := sum / float64(len(scores))
	pow := math.Pow10(precision)
	result := math.Round(avg*pow) / pow
	return result
}

// midPoint eliminates time gaps between the end of one verse and the beginning of the next.
func (m *MMSAlign) midPoint(parts []db.Audio) []db.Audio {
	for i := range parts {
		if i == 0 {
			parts[0].BeginTS = 0.0
		} else {
			midPoint := (parts[i].BeginTS + parts[i-1].EndTS) / 2.0
			parts[i].BeginTS = midPoint
			parts[i-1].EndTS = midPoint
		}
	}
	return parts
}
