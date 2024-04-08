package speech_to_text

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Docs:
https://github.com/openai/whisper
Install:
pip3 install git+https://github.com/openai/whisper.git
Whisper is an open source Speech to Text program developed by OpenAI.
Executable:
/Users/gary/Library/Python/3.9/bin/whisper
*/

type Whisper struct {
	ctx       context.Context
	conn      db.DBAdapter
	bibleId   string
	model     string
	outputDir string
	records   []db.Script
}

func NewWhisper(bibleId string, conn db.DBAdapter, model string) Whisper {
	var w Whisper
	w.ctx = conn.Ctx
	w.conn = conn
	w.bibleId = bibleId
	w.model = model
	w.records = make([]db.Script, 0, 100000)
	return w
}

func (w *Whisper) ProcessDirectory(filesetId string, testament dataset.TestamentType) dataset.Status {
	var status dataset.Status
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), w.bibleId, filesetId)
	w.outputDir = directory + `_whisper`
	files, err := os.ReadDir(directory)
	if err != nil {
		return log.Error(w.ctx, 500, err, `Error reading directory`)
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, `.`) {
			fmt.Println(filename)
			fileType := filename[:1]
			if fileType == `A` && (testament == dataset.OT || testament == dataset.C) {
				w.processFile(directory, filename)
			} else if fileType == `B` && (testament == dataset.NT || testament == dataset.C) {
				w.processFile(directory, filename)
			}
		}
	}
	w.loadWhisperOutput(w.outputDir)
	return status
}

func (w *Whisper) processFile(directory string, filename string) dataset.Status {
	bookId, chapter, status := w.parseFilename(filename)
	if status.IsErr {
		return status
	}
	if bookId == `TIT` && chapter == 3 {
		var path = filepath.Join(directory, filename)
		w.runWhisper(path)
	}
	return status
}

func (w *Whisper) runWhisper(audioFilePath string) dataset.Status {
	var status dataset.Status
	whisperPath := os.Getenv(`WHISPER_EXE`)
	cmd := exec.Command(whisperPath, audioFilePath,
		`--model`, w.model,
		`--output_format`, `json`,
		`--output_dir`, w.outputDir)
	// --language is another option
	fmt.Println(cmd.String())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		return log.Error(w.ctx, 500, err, `Error running Whisper`)
	}
	stderrStr := stderrBuf.String()
	if stderrStr != `` {
		log.Warn(w.ctx, err, `Stderr message when running Whisper`)
	}
	return status
}

func (w *Whisper) loadWhisperOutput(directory string) dataset.Status {
	var status dataset.Status
	type WhisperSegmentType struct {
		Id     int     `json:"id"`
		Seek   float64 `json:"seek"`
		Start  float64 `json:"start"`
		End    float64 `json:"end"`
		Text   string  `json:"text"`
		Tokens []int   `json:"tokens"`
		// "tokens": [50660, 293, 281, 12076, 11, 281, 312, 42541, 11, 281, 312, 1919, 337, 633, 665, 589, 11, 281, 1710, 6724, 295, 51004],
		Temperature      float32 `json:"temperature"`
		AvgLogProb       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	}
	type WhisperOutputType struct {
		Segments []WhisperSegmentType `json:"segments"`
		Language string               `json:"language"`
	}
	jsonFiles, err := os.ReadDir(directory)
	if err != nil {
		return log.Error(w.ctx, 500, err, `Error reading directory`)
	}
	for _, jsonFile := range jsonFiles {
		filePath := filepath.Join(directory, jsonFile.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return log.Error(w.ctx, 500, err, `Error reading file`)
		}
		var response WhisperOutputType
		err = json.Unmarshal(content, &response)
		if err != nil {
			return log.Error(w.ctx, 500, err, "Error decoding Whisper JSON")
		}
		for i, seg := range response.Segments {
			var rec db.Script
			rec.BookId, rec.ChapterNum, status = w.parseFilename(jsonFile.Name())
			rec.AudioFile = jsonFile.Name()
			rec.ScriptNum = strconv.Itoa(i + 1) // Works because it process 1 chapter per file.
			rec.VerseNum = 0
			rec.VerseStr = `0`
			rec.ScriptTexts = []string{seg.Text}
			rec.ScriptBeginTS = seg.Start
			rec.ScriptEndTS = seg.End
			w.records = append(w.records, rec)
		}
	}
	w.conn.InsertScripts(w.records)
	return status
}

func (w *Whisper) parseFilename(filename string) (string, int, dataset.Status) {
	var bookId string
	var chapter int
	var status dataset.Status
	chapter, err := strconv.Atoi(filename[6:8])
	if err != nil {
		status = log.Error(w.ctx, 500, err, `Error parsing chapter num`)
		return bookId, chapter, status
	}
	bookName := strings.Replace(filename[9:21], `_`, ``, -1)
	bookId = db.USFMBookId(w.ctx, bookName)
	return bookId, chapter, status
}
