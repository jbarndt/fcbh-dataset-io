package speech_to_text

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	ctx     context.Context
	conn    db.DBAdapter
	bibleId string
	model   string
}

func NewWhisper(bibleId string, conn db.DBAdapter, model string) Whisper {
	var w Whisper
	w.ctx = conn.Ctx
	w.conn = conn
	w.bibleId = bibleId
	w.model = model
	return w
}

func (w *Whisper) ProcessFiles(files []input.InputFile) dataset.Status {
	var status dataset.Status
	var outputFile string
	for _, file := range files {
		outputFile, status = w.RunWhisper(file)
		w.loadWhisperOutput(outputFile, file)
	}
	return status
}

func (w *Whisper) RunWhisper(audioFile input.InputFile) (string, dataset.Status) {
	var outputDir, status = w.ensureOutputDir(audioFile)
	if status.IsErr {
		return outputDir, status
	}
	whisperPath := os.Getenv(`WHISPER_EXE`)
	cmd := exec.Command(whisperPath, audioFile.FilePath(),
		`--model`, w.model,
		`--output_format`, `json`,
		`--output_dir`, outputDir)
	// --language is another option
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(w.ctx, 500, err, `Error running Whisper`)
		// Do not return immediately, must get std error
	}
	stderrStr := stderrBuf.String()
	if stderrStr != `` {
		log.Warn(w.ctx, `Whisper Stderr:`, stderrStr)
	}
	fileType := filepath.Ext(audioFile.Filename)
	outputFile := filepath.Join(outputDir, audioFile.Filename[:len(audioFile.Filename)-len(fileType)]) + `.json`
	return outputFile, status
}

func (w *Whisper) ensureOutputDir(audioFile input.InputFile) (string, dataset.Status) {
	var status dataset.Status
	var outputDir = audioFile.Directory + `_WHISPER`
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0777)
	} else if err != nil {
		status = log.Error(w.ctx, 500, err, `Error creating whisper output directory`)
	}
	return outputDir, status
}

func (w *Whisper) loadWhisperOutput(outputFile string, file input.InputFile) dataset.Status {
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
	var records = make([]db.Script, 0, 100)
	content, err := os.ReadFile(outputFile)
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
		rec.BookId = file.BookId
		rec.ChapterNum = file.Chapter
		rec.AudioFile = file.Filename
		rec.ScriptNum = strconv.Itoa(i + 1) // Works because it process 1 chapter per file.
		rec.VerseNum = 0
		rec.VerseStr = `0`
		rec.ScriptTexts = []string{seg.Text}
		rec.ScriptBeginTS = seg.Start
		rec.ScriptEndTS = seg.End
		records = append(records, rec)
	}
	w.conn.InsertScripts(records)
	return status
}
