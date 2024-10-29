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
	/*
		Sandeep timestamp solution
		// Get Timestamps from Sandeep's bucket
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
	*/
	// Waha Solution
	bucket := timestamp.NewWahaTimestamper(a.ctx)
	timestamps, status = bucket.GetTimestamps(timestamp.VerseAeneas, file.MediaId, file.BookId, file.Chapter)
	// My FA Timestamps
	//bucket := timestamp.NewFATimeStamper(a.ctx, a.conn)
	//timestamps, status = bucket.GetTimestamps(file.BookId, file.Chapter)
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
		fmt.Println(ts.BookId, ts.ChapterNum, ts.VerseStr, response)
		timestamps[i].Text = response
	}
	//a.conn.InsertTimestamps(timestamps)
	bytes, err := json.Marshal(timestamps)
	if err != nil {
		return log.Error(a.ctx, 500, err, `Error marshalling timestamps`)
	}
	fmt.Println(file.BookId, file.Chapter, string(bytes))
	return status
}
