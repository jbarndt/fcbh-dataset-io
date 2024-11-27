package mms

import (
	"bufio"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"dataset/timestamp"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MMSFA2_Input struct {
	AudioFile string   `json:"audio_file"`
	NormWords []string `json:"words"`
}

type Word struct {
	scriptId   int64
	wordSeq    int
	word       string
	uroman     string
	normalized string
}

type Timestamp struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Score float64 `json:"score"`
}

type MMSFA2 struct {
	ctx     context.Context
	conn    db.DBAdapter // This database adapter must contain the text to be processed
	lang    string
	sttLang string // I don't know if this is useful
}

func NewMMSFA2(ctx context.Context, conn db.DBAdapter, lang string, sttLang string) MMSFA2 {
	var m MMSFA2
	m.ctx = ctx
	m.conn = conn
	m.lang = lang
	m.sttLang = sttLang
	return m
}

// ProcessFiles will perform Forced Alignment on these files
func (a *MMSFA2) ProcessFiles(files []input.InputFile) dataset.Status {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_fa2.py")
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
func (m *MMSFA2) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) dataset.Status {
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_fa_")
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var faInput MMSFA2_Input
	faInput.AudioFile, status = timestamp.ConvertMp3ToWav(m.ctx, tempDir, file.FilePath())
	if status.IsErr {
		return status
	}
	var verses []db.Script
	verses, status = m.conn.SelectScriptsByChapter(file.BookId, file.Chapter)
	if status.IsErr {
		return status
	}
	var wordList []Word
	faInput.NormWords, wordList, status = m.prepareText(m.lang, verses)
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
		return log.Error(m.ctx, 500, err, "Error writing to mms_fa.py")
	}
	err = writer.Flush()
	if err != nil {
		return log.Error(m.ctx, 500, err, "Error flush to mms_fa.py")
	}
	response, err2 := reader.ReadString('\n')
	if err2 != nil {
		return log.Error(m.ctx, 500, err2, `Error reading mms_fa.py response`)
	}
	m.processPyOutput(file, wordList, response)
	// development
	//err = os.WriteFile("engweb_fa_out.json", []byte(response), 0644)
	//if err != nil {
	//	panic(err)
	//}
	return status
}

func (m *MMSFA2) prepareText(lang string, scripts []db.Script) ([]string, []Word, dataset.Status) {
	var textList []string
	var refList []Word
	var status dataset.Status
	re1 := regexp.MustCompile(`[^a-z' ]`)
	re2 := regexp.MustCompile(` +`)
	var verses, uroman []string
	for _, script := range scripts {
		verses = append(verses, script.ScriptText)
	}
	uroman, status = URoman(m.ctx, lang, verses)
	for i, text := range uroman {
		text = strings.ToLower(text)
		text = strings.ReplaceAll(text, "'", "'")
		text = re1.ReplaceAllString(text, " ")
		text = re2.ReplaceAllString(text, " ")
		text = strings.TrimSpace(text)
		norm := strings.Fields(text)
		words := strings.Fields(scripts[i].ScriptText)
		urom := strings.Fields(uroman[i])
		if len(words) != len(norm) {
			status = log.ErrorNoErr(m.ctx, 500, "Word count did not match in MMS_FA prepareText", len(words), len(norm))
		}
		if len(words) != len(urom) {
			status = log.ErrorNoErr(m.ctx, 500, "Uroman count did not match in MMS_FA prepareText", len(words), len(urom))
		}
		for w := range words {
			textList = append(textList, norm[w])
			var ref Word
			ref.scriptId = int64(scripts[i].ScriptId)
			ref.wordSeq = w
			ref.word = words[w]
			ref.uroman = urom[w]
			ref.normalized = norm[w]
			refList = append(refList, ref)
		}
	}
	return textList, refList, status
}

func (m *MMSFA2) processPyOutput(file input.InputFile, wordRefs []Word, response string) dataset.Status {
	var status dataset.Status
	response = strings.TrimRight(response, "\n")
	var timestamps []Timestamp
	err := json.Unmarshal([]byte(response), &timestamps)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error unmarshalling json`)
	}
	if len(timestamps) != len(wordRefs) {
		return log.ErrorNoErr(m.ctx, 400, "Num words input to mms_fs:", len(wordRefs), ", num timestamps returned:", len(timestamps))
	}
	var words []db.Audio
	for i, ref := range wordRefs {
		var word db.Audio
		word.BookId = file.BookId
		word.ChapterNum = file.Chapter
		word.AudioFile = file.Filename
		word.ScriptId = ref.scriptId
		//word.VerseSeq =
		word.WordSeq = ref.wordSeq
		word.Text = ref.word
		word.Uroman = ref.uroman
		//ref.normalized
		word.BeginTS = timestamps[i].Start
		word.EndTS = timestamps[i].End
		word.FAScore = timestamps[i].Score
		words = append(words, word)
	}
	var wordsByVerse [][]db.Audio
	wordsByVerse = m.groupByVerse(words)
	var verses []db.Audio
	verses = m.summarizeByVerse(wordsByVerse)
	verses = m.addSpace(verses)
	status = m.conn.UpdateScriptFATimestamps(verses)
	if status.IsErr {
		return status
	}
	words, status = m.conn.InsertAudioWords(words)
	return status
}

func (m *MMSFA2) groupByVerse(words []db.Audio) [][]db.Audio {
	var result [][]db.Audio
	var verse []db.Audio
	var verseSeq = 0
	for i, word := range words {
		if word.WordSeq == 0 {
			if i > 0 {
				result = append(result, verse)
				verse = nil
				verseSeq++
			}
		}
		word.VerseSeq = verseSeq
		verse = append(verse, word)
	}
	result = append(result, verse)
	return result
}

func (m *MMSFA2) summarizeByVerse(chapter [][]db.Audio) []db.Audio {
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
		vs.FAScore = m.average(scores, 3)
		result = append(result, vs)
	}
	return result
}

func (m *MMSFA2) average(scores []float64, precision int) float64 {
	var sum float64
	for _, scr := range scores {
		sum += scr
	}
	avg := sum / float64(len(scores))
	pow := math.Pow10(precision)
	result := math.Round(avg*pow) / pow
	return result
}

// addSpace eliminates time gaps between the end of one verse and the beginning of the next.
func (m *MMSFA2) addSpace(parts []db.Audio) []db.Audio {
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
