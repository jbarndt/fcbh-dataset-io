package speech_to_text

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"encoding/json"
	"fmt"
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

type WhisperVs struct {
	ctx     context.Context
	conn    db.DBAdapter
	bibleId string
	model   string
	lang2   string // 2 char language code
	tempDir string
}

func NewWhisperVs(bibleId string, conn db.DBAdapter, model string, lang2 string) WhisperVs {
	var w WhisperVs
	w.ctx = conn.Ctx
	w.conn = conn
	w.bibleId = bibleId
	w.model = model
	w.lang2 = lang2
	return w
}

func (w *WhisperVs) ProcessFiles(files []input.InputFile) dataset.Status {
	var status dataset.Status
	var outputFile string
	var err error
	w.tempDir, err = os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "WhisperVs_")
	if err != nil {
		return log.Error(w.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(w.tempDir)
	for _, file := range files {
		fmt.Println(`INPUT FILE:`, file)
		var pieces []db.Timestamp
		pieces, status = w.CopyByTimestamp(file)
		if status.IsErr {
			return status
		}
		status = w.conn.DeleteScripts(file.BookId, file.Chapter)
		if status.IsErr {
			return status
		}
		var records []db.Script
		for pieceNum, piece := range pieces {
			fmt.Println(`VERSE PIECE:`, piece)
			outputFile, status = w.RunWhisper(piece)
			var rec db.Script
			rec, status = w.loadWhisperOutput(outputFile, file, pieceNum, piece)
			records = append(records, rec)
		}
		w.conn.InsertScripts(records)
		records = nil
	}
	return status
}

func (w *WhisperVs) CopyByTimestamp(file input.InputFile) ([]db.Timestamp, dataset.Status) {
	var results []db.Timestamp
	var status dataset.Status
	timestamps, status := w.conn.SelectScriptTimestamps(file.BookId, file.Chapter)
	if status.IsErr {
		return results, status
	}
	if len(timestamps) == 0 {
		status = log.ErrorNoErr(w.ctx, 400, `No timestamps found for`, file.BookId, file.Chapter)
		return results, status
	}
	var command []string
	command = append(command, `-i`, file.FilePath())
	command = append(command, `-codec:a`, `copy`)
	command = append(command, `-y`)
	for _, ts := range timestamps {
		if ts.BeginTS == 0.0 && ts.EndTS == 0.0 {
			continue
		}
		beginTS := strconv.FormatFloat(ts.BeginTS, 'f', 2, 64)
		command = append(command, `-ss`, beginTS)
		if ts.EndTS != 0.0 {
			endTS := strconv.FormatFloat(ts.EndTS, 'f', 2, 64)
			command = append(command, `-to`, endTS)
		}
		ts.AudioFile = fmt.Sprintf("verse_%s_%d_%s_%s.mp3",
			file.BookId, file.Chapter, ts.VerseStr, beginTS)
		outputPath := filepath.Join(w.tempDir, ts.AudioFile)
		command = append(command, `-c`, `copy`, outputPath)
		results = append(results, ts)
	}
	ffMpegPath := `ffmpeg`
	cmd := exec.Command(ffMpegPath, command...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(w.ctx, 500, err, stderrBuf.String())
	}
	return results, status
}

func (w *WhisperVs) RunWhisper(audio db.Timestamp) (string, dataset.Status) {
	var status dataset.Status
	whisperPath := os.Getenv(`WHISPER_EXE`)
	cmd := exec.Command(whisperPath,
		filepath.Join(w.tempDir, audio.AudioFile),
		`--model`, w.model,
		`--output_format`, `json`,
		`--fp16`, `False`,
		`--language`, w.lang2,
		`--output_dir`, w.tempDir)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(w.ctx, 500, err, stderrBuf.String())
		// Do not return immediately, must get std error
	}
	stderrStr := stderrBuf.String()
	if stderrStr != `` {
		log.Warn(w.ctx, `Whisper Stderr:`, stderrStr)
	}
	fileType := filepath.Ext(audio.AudioFile)
	outputFile := audio.AudioFile[:len(audio.AudioFile)-len(fileType)] + `.json`
	return outputFile, status
}

func (w *WhisperVs) loadWhisperOutput(outputFile string, file input.InputFile,
	pieceNum int, piece db.Timestamp) (db.Script, dataset.Status) {
	var rec db.Script
	var status dataset.Status
	type WhisperSegmentType struct {
		Id               int     `json:"id"`
		Seek             float64 `json:"seek"`
		Start            float64 `json:"start"`
		End              float64 `json:"end"`
		Text             string  `json:"text"`
		Tokens           []int   `json:"tokens"`
		Temperature      float32 `json:"temperature"`
		AvgLogProb       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	}
	type WhisperOutputType struct {
		Segments []WhisperSegmentType `json:"segments"`
		Language string               `json:"language"`
	}
	content, err := os.ReadFile(filepath.Join(w.tempDir, outputFile))
	if err != nil {
		return rec, log.Error(w.ctx, 500, err, `Error reading file`)
	}
	var response WhisperOutputType
	err = json.Unmarshal(content, &response)
	if err != nil {
		return rec, log.Error(w.ctx, 500, err, "Error decoding Whisper JSON")
	}
	rec.BookId = file.BookId
	rec.ChapterNum = file.Chapter
	rec.AudioFile = file.Filename
	rec.ScriptNum = strconv.Itoa(pieceNum + 1)
	rec.VerseNum = dataset.SafeVerseNum(piece.VerseStr)
	rec.VerseStr = piece.VerseStr
	for i, seg := range response.Segments {
		rec.ScriptTexts = append(rec.ScriptTexts, seg.Text)
		if i == 0 {
			rec.ScriptBeginTS = seg.Start + piece.BeginTS
		}
		rec.ScriptEndTS = seg.End + piece.BeginTS
	}
	return rec, status
}
