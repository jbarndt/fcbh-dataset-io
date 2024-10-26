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
	"fmt"
	"os"
	"strings"
)

type MMSFA_Input struct {
	BookId    string         `json:"book_id"`
	Chapter   int            `json:"chapter"`
	AudioFile string         `json:"audio_file"`
	Verses    []MMSFA_Verses `json:"verses"`
}

type MMSFA_Verses struct {
	Verse string `json:"verse_str"`
	Text  string `json:"text"`
}

type MMSFA_Output struct {
	BookId  string  `json:"book"`
	Chapter int     `json:"chapter"`
	Verse   string  `json:"verse"`
	Start   float64 `json:"start"`
	End     float64 `json:"end"`
	Score   float64 `json:"score"`
	Text    string  `json:"text"`
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
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang) // is this correct for mms_fa
	if status.IsErr {
		return status
	}
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_PYTHON`), "mms_fa.py", lang)
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

// processFile
func (m *MMSFA) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) dataset.Status {
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_fa_")
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var faInput MMSFA_Input
	faInput.BookId = file.BookId
	faInput.Chapter = file.Chapter
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
		faVerse.Verse = vers.VerseStr
		faVerse.Text = vers.ScriptText
		faInput.Verses = append(faInput.Verses, faVerse)
	}
	content, err := json.Marshal(faInput)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error marshalling json`)
	}
	err2 := os.WriteFile("engweb_fa_test.json", content, 0644)
	if err2 != nil {
		panic(err2)
	}
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
	response = strings.TrimRight(response, "\n")
	var output MMSFA_Output
	err = json.Unmarshal([]byte(response), &output)
	if err != nil {
		return log.Error(m.ctx, 500, err, `Error unmarshalling json`)
	}
	fmt.Println(output)
	return status
}
