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
	"strings"
)

type MMSFA_Input struct {
	AudioFile string         `json:"audio_file"`
	Verses    []MMSFA_Verses `json:"verses"`
}

type MMSFA_Verses struct {
	ScriptId int    `json:"script_id"`
	Text     string `json:"text"`
}

type MMSFA struct {
	ctx     context.Context
	conn    db.DBAdapter // This database adapter must contain the text to be processed
	lang    string
	sttLang string // I don't know if this is useful
}

func NewMMSFA(ctx context.Context, conn db.DBAdapter, lang string, sttLang string) MMSFA {
	var m MMSFA
	m.ctx = ctx
	m.conn = conn
	m.lang = lang
	m.sttLang = sttLang
	return m
}

// ProcessFiles will perform Forced Alignment on these files
func (a *MMSFA) ProcessFiles(files []input.InputFile) dataset.Status {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_fa.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_FA_PYTHON`), pythonScript, lang)
	if status.IsErr {
		return status
	}
	for _, file := range files {
		status = a.processFile(file, writer, reader)
		if status.IsErr {
			return status
		}
	}
	return status
}

// processFile will process one audio file through mms forced alignment
func (m *MMSFA) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) dataset.Status {
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_fa_")
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var faInput MMSFA_Input
	faInput.AudioFile, status = timestamp.ConvertMp3ToWav(m.ctx, tempDir, file.FilePath())
	if status.IsErr {
		return status
	}
	var verses []db.Script
	verses, status = m.conn.SelectScriptsByChapter(file.BookId, file.Chapter)
	if status.IsErr {
		return status
	}
	for _, vers := range verses {
		var faVerse MMSFA_Verses
		faVerse.ScriptId = vers.ScriptId
		faVerse.Text = vers.ScriptText
		faInput.Verses = append(faInput.Verses, faVerse)
	}
	content, err := json.Marshal(faInput)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error marshalling json`)
	}
	// temp
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
	m.processPyOutput(file, response)
	return status
}

func (m *MMSFA) processPyOutput(file input.InputFile, response string) dataset.Status {
	var status dataset.Status
	response = strings.TrimRight(response, "\n")
	var words []db.Audio
	err := json.Unmarshal([]byte(response), &words)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error unmarshalling json`)
	}
	wordsByVerse := m.groupByVerse(words)
	verses := m.summarizeByVerse(wordsByVerse)
	for i := range verses {
		verses[i].AudioFile = file.Filename
	}
	verses = m.addSpace(verses)
	status = m.conn.UpdateScriptFATimestamps(verses)
	if status.IsErr {
		return status
	}
	for i := range verses {
		_, status = m.conn.InsertAudioWords(verses[i], wordsByVerse[i])
		if status.IsErr {
			return status
		}
	}
	return status
}

func (m *MMSFA) groupByVerse(words []db.Audio) [][]db.Audio {
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

func (m *MMSFA) summarizeByVerse(chapter [][]db.Audio) []db.Audio {
	var result []db.Audio
	for _, verse := range chapter {
		var vs = verse[0]
		vs.EndTS = verse[len(verse)-1].EndTS
		var scores []float64
		var uroman []string
		for _, word := range verse {
			scores = append(scores, word.FAScore)
			uroman = append(uroman, word.Uroman)
		}
		vs.FAScore = m.average(scores, 3)
		vs.Uroman = strings.Join(uroman, " ")
		result = append(result, vs)
	}
	return result
}

func (m *MMSFA) average(scores []float64, precision int) float64 {
	var sum float64
	for _, scr := range scores {
		sum += scr
	}
	avg := sum / float64(len(scores))
	pow := math.Pow10(precision)
	result := math.Round(avg*pow) / pow
	return result
}

// addSpace adds back space that was removed from words by mms_fa
func (m *MMSFA) addSpace(parts []db.Audio) []db.Audio {
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
