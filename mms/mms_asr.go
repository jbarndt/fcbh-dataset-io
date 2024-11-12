package mms

import (
	"bufio"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"dataset/timestamp"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MMSASR struct {
	ctx     context.Context
	conn    db.DBAdapter
	lang    string
	sttLang string
}

func NewMMSASR(ctx context.Context, conn db.DBAdapter, lang string, sttLang string) MMSASR {
	var a MMSASR
	a.ctx = ctx
	a.conn = conn
	a.lang = lang
	a.sttLang = sttLang
	return a
}

// ProcessFiles will perform Auto Speech Recognition on these files
func (a *MMSASR) ProcessFiles(files []input.InputFile) dataset.Status {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_asr.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_ASR_PYTHON`), pythonScript, lang)
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
func (a *MMSASR) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) dataset.Status {
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_asr_")
	if err != nil {
		return log.Error(a.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	wavFile, status := timestamp.ConvertMp3ToWav(a.ctx, tempDir, file.FilePath())
	if status.IsErr {
		return status
	}
	var timestamps []db.Audio
	timestamps, status = a.conn.SelectFAScriptTimestamps(file.BookId, file.Chapter)
	if status.IsErr {
		return status
	}
	timestamps, status = timestamp.ChopByTimestamp(a.ctx, tempDir, wavFile, timestamps)
	for i, ts := range timestamps {
		timestamps[i].AudioFile = file.Filename
		timestamps[i].AudioChapterWav = wavFile
		_, err = writer.WriteString(ts.AudioVerseWav + "\n")
		if err != nil {
			return log.Error(a.ctx, 500, err, "Error writing to mms_asr.py")
		}
		err = writer.Flush()
		if err != nil {
			return log.Error(a.ctx, 500, err, "Error flush to mms_asr.py")
		}
		response, err2 := reader.ReadString('\n')
		if err2 != nil {
			return log.Error(a.ctx, 500, err2, `Error reading mms_asr.py response`)
		}
		response = strings.TrimRight(response, "\n")
		fmt.Println(ts.BookId, ts.ChapterNum, ts.VerseStr, ts.ScriptId, response)
		timestamps[i].Text = response
	}
	var recCount int
	recCount, status = a.conn.UpdateScriptText(timestamps)
	fmt.Println(file.BookId, file.Chapter, recCount)
	return status
}
