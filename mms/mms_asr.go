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
	"github.com/garygriswold/lang_tree/search"
	"io"
	"os"
	"os/exec"
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
	// when entirely done, send it an exit message, such as ctrl-D
	return status
}

// Check that language is supported by mms_asr, and return alternate if it is not
// This method should be in mms_util.go
func checkLanguage(ctx context.Context, lang string, sttLang string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	if sttLang != `` {
		result = sttLang
	} else {
		var tree = search.NewLanguageTree(ctx)
		err := tree.Load()
		if err != nil {
			status = log.Error(ctx, 500, err, `Error loading language`)
			return result, status
		}
		langs, distance, err2 := tree.Search(strings.ToLower(lang), "mms_asr")
		if err2 != nil {
			status = log.Error(ctx, 500, err2, `Error Searching for language`)
		}
		if len(langs) > 0 {
			result = langs[0]
			log.Info(ctx, `Using language`, result, "distance:", distance)
		} else {
			status = log.ErrorNoErr(ctx, 400, `No compatible language code was found for`, lang)
		}
	}
	return result, status
}

// callPythonScript will exec the python script, and setup pipes on stdin and stdout
// this method should be moved to util
func callStdIOScript(ctx context.Context, command string, arg ...string) (io.Writer, io.Reader, dataset.Status) {
	var writer io.Writer
	var reader io.Reader
	var status dataset.Status
	cmd := exec.Command(command, arg...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to open stdin for writing to Fasttext`)
		return writer, reader, status
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to open stdout for writing to Fasttext`)
		return writer, reader, status
	}
	err = cmd.Start()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to start writing to Fasttext`)
		return writer, reader, status
	}
	writer = bufio.NewWriterSize(stdin, 4096)
	reader = bufio.NewReaderSize(stdout, 4096)
	return writer, reader, status
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
	var timestamps []timestamp.Timestamp
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
