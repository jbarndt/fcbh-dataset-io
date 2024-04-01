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
	bibleId   string
	conn      db.DBAdapter
	outputDir string
	records   []db.Script
}

func NewWhisper(bibleId string, conn db.DBAdapter) Whisper {
	var w Whisper
	w.ctx = conn.Ctx
	w.bibleId = bibleId
	w.conn = conn
	w.records = make([]db.Script, 0, 100000)
	return w
}

func (w *Whisper) ProcessDirectory(filesetId string, testament dataset.TestamentType) {
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), w.bibleId, filesetId)
	w.outputDir = directory + `_whisper`
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Error(w.ctx, err)
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
}

func (w *Whisper) processFile(directory string, filename string) {
	bookId, chapter := w.parseFilename(filename)
	if bookId == `TIT` && chapter == 3 {
		var path = filepath.Join(directory, filename)
		w.runWhisper(path)
	}
}

func (w *Whisper) runWhisper(audioFilePath string) {
	whisperPath := os.Getenv(`WHISPER_EXE`)
	cmd := exec.Command(whisperPath, audioFilePath,
		`--model`, `tiny`,
		`--output_format`, `json`,
		`--output_dir`, w.outputDir)
	// --language is another option
	fmt.Println(cmd.String())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	stderrStr := stderrBuf.String()
	if stderrStr != `` {
		fmt.Printf("Stderr: \n%s\n", stderrStr)
	}
}

func (w *Whisper) loadWhisperOutput(directory string) {
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
		log.Error(w.ctx, err)
	}
	for _, jsonFile := range jsonFiles {
		filePath := filepath.Join(directory, jsonFile.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Error(w.ctx, err)
		}
		var response WhisperOutputType
		err = json.Unmarshal(content, &response)
		if err != nil {
			log.Error(w.ctx, "Error decoding Whisper JSON:", err)
		}
		for i, seg := range response.Segments {
			//fmt.Println(seg)
			var rec db.Script
			rec.BookId, rec.ChapterNum = w.parseFilename(jsonFile.Name())
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
}

func (w *Whisper) parseFilename(filename string) (string, int) {
	chapter, err := strconv.Atoi(filename[6:8])
	if err != nil {
		log.Fatal(w.ctx, err)
	}
	bookName := strings.Replace(filename[9:21], `_`, ``, -1)
	bookId := db.USFMBookId(w.ctx, bookName)
	return bookId, chapter
}
