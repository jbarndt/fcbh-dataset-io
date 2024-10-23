package mms

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"dataset/timestamp"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang)
	if status.IsErr {
		return status
	}
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_PYTHON`), "mms_asr.py", lang)
	if status.IsErr {
		return status
	}
	for _, file := range files {
		status = a.processFile(file, writer, reader)
		if status.IsErr {
			return status
		}
	}
	// when entirely done, send it an exit message, such as ctrl-D
	return status
}

// processFile
func (a *MMSASR) processFile(file input.InputFile, writer io.Writer, reader io.Reader) dataset.Status {
	//var scripts []db.Script
	var status dataset.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_asr_")
	if err != nil {
		return log.Error(a.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var bucket timestamp.TSBucket
	bucket, status = timestamp.NewTSBucket(a.ctx)
	if status.IsErr {
		return status
	}
	var timestamps []db.Audio
	timestamps, status = bucket.GetTimestamps(timestamp.VerseAeneas, file.MediaId, file.BookId, file.Chapter)
	if status.IsErr {
		return status
	}
	timestamps, status = timestamp.ChopByTimestamp(a.ctx, tempDir, file, timestamps)
	for i, ts := range timestamps {
		writer.Write([]byte(ts.AudioVerse))
		var response []byte
		count, err2 := reader.Read(response)
		if err2 != nil {
			return log.Error(a.ctx, 500, err2, `Error reading response`)
		}
		fmt.Println("count:", count)
		timestamps[i].Text = string(response[:count])
	}
	//a.conn.InsertTimestamps(timestamps)
	bytes, err := json.Marshal(timestamps)
	if err != nil {
		return log.Error(a.ctx, 500, err, `Error marshalling timestamps`)
	}
	fmt.Println(file.BookId, file.Chapter, string(bytes))
	return status
}
