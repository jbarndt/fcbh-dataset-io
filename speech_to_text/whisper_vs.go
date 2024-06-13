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
}

func NewWhisperVs(bibleId string, conn db.DBAdapter, model string) WhisperVs {
	var w WhisperVs
	w.ctx = conn.Ctx
	w.conn = conn
	w.bibleId = bibleId
	w.model = model
	return w
}

func (w *WhisperVs) ProcessFiles(files []input.InputFile) dataset.Status {
	var status dataset.Status
	var outputFile string
	for _, file := range files {
		fmt.Println(`INPUT FILE:`, file)
		var pieces []db.Timestamp
		pieces, status = w.ChopByTimestamp(file)
		if status.IsErr {
			return status
		}
		status = w.conn.DeleteScripts(file.BookId, file.Chapter)
		if status.IsErr {
			return status
		}
		for pieceNum, piece := range pieces {
			fmt.Println(`VERSE PIECE:`, piece)
			outputFile, status = w.RunWhisper(piece)
			status = w.loadWhisperOutput(outputFile, file, pieceNum, piece)
		}
	}
	return status
}

func (w *WhisperVs) ChopByTimestamp(audioFile input.InputFile) ([]db.Timestamp, dataset.Status) {
	var results []db.Timestamp
	var status dataset.Status
	timestamps, status := w.conn.SelectScriptTimestamps(audioFile.BookId, audioFile.Chapter)
	if status.IsErr {
		return results, status
	}
	ffMpegPath := `ffmpeg`
	for _, ts := range timestamps {
		if ts.BeginTS == 0.0 && ts.EndTS == 0.0 {
			continue
		}
		var cmd *exec.Cmd
		beginTS := strconv.FormatFloat(ts.BeginTS, 'g', -1, 64)
		endTS := strconv.FormatFloat(ts.EndTS, 'g', -1, 64)

		outputFile := fmt.Sprintf("output_%v_%v.mp3", ts.BeginTS, ts.EndTS)
		ts.AudioFile = filepath.Join(os.Getenv(`FCBH_DATASET_TMP`), outputFile)
		if ts.EndTS != 0.0 {
			cmd = exec.Command(ffMpegPath, `-y`, `-i`, audioFile.FilePath(),
				`-ss`, beginTS, `-to`, endTS, `-c`, `copy`, ts.AudioFile)
		} else {
			cmd = exec.Command(ffMpegPath, `-y`, `-i`, audioFile.FilePath(),
				`-ss`, beginTS, `-c`, `copy`, ts.AudioFile)
		}
		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
		err := cmd.Run()
		if err != nil {
			status = log.Error(w.ctx, 500, err, stderrBuf.String())
			return results, status
		}
		results = append(results, ts)
	}
	return results, status
}

func (w *WhisperVs) RunWhisper(audio db.Timestamp) (string, dataset.Status) {
	var status dataset.Status
	outputDir := os.Getenv(`FCBH_DATASET_TMP`)
	//var outputDir, status = w.ensureOutputDir(audio.AudioFile)
	//if status.IsErr {
	//	return outputDir, status
	//}
	whisperPath := os.Getenv(`WHISPER_EXE`)
	cmd := exec.Command(whisperPath, audio.AudioFile,
		`--model`, w.model,
		`--output_format`, `json`,
		`--fp16`, `False`,
		`--output_dir`, outputDir)
	// --language is another option
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

// Should this be a user directory under tmp??
func (w *WhisperVs) ensureOutputDir(audioFile string) (string, dataset.Status) {
	var status dataset.Status
	var outputDir = filepath.Dir(audioFile) + `_WHISPER`
	//var outputDir = audioFile.Directory + `_WHISPER`
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0777)
	} else if err != nil {
		status = log.Error(w.ctx, 500, err, `Error creating whisper output directory`)
	}
	return outputDir, status
}

func (w *WhisperVs) loadWhisperOutput(outputFile string, file input.InputFile,
	pieceNum int, piece db.Timestamp) dataset.Status {
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
	var rec db.Script
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
	records = append(records, rec)
	status = w.conn.InsertScripts(records)
	return status
}
