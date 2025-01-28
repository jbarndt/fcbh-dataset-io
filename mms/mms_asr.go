package mms

import (
	"bufio"
	"context"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"dataset/timestamp"
	"dataset/utility"
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
	uroman  utility.StdioExec
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
func (a *MMSASR) ProcessFiles(files []input.InputFile) *log.Status {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status != nil {
		return status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_asr.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_ASR_PYTHON`), pythonScript, lang)
	if status != nil {
		return status
	}
	uromanPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "mms", "uroman_stdio.py")
	a.uroman, status = utility.NewStdioExec(a.ctx, os.Getenv(`FCBH_MMS_FA_PYTHON`), uromanPath, "-l", a.lang)
	if status != nil {
		return status
	}
	defer a.uroman.Close()
	for _, file := range files {
		log.Info(a.ctx, "MMS ASR", file.BookId, file.Chapter)
		status = a.processFile(file, writer, reader)
		if status != nil {
			return status
		}
	}
	return status
}

// processFile
func (a *MMSASR) processFile(file input.InputFile, writer *bufio.Writer, reader *bufio.Reader) *log.Status {
	var status *log.Status
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_asr_")
	if err != nil {
		return log.Error(a.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	wavFile, status := timestamp.ConvertMp3ToWav(a.ctx, tempDir, file.FilePath())
	if status != nil {
		return status
	}
	var timestamps []db.Audio
	timestamps, status = a.conn.SelectFAScriptTimestamps(file.BookId, file.Chapter)
	if status != nil {
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
		uRoman, status2 := a.uroman.Process(response)
		if status2 != nil {
			return status2
		}
		timestamps[i].Uroman = uRoman
	}
	//log.Debug(a.ctx, "Finished ASR", file.BookId, file.Chapter)
	var recCount int
	recCount, status = a.conn.UpdateScriptText(timestamps)
	if recCount != len(timestamps) {
		log.Warn(a.ctx, "Timestamp update counts need investigation", recCount, len(timestamps))
	}
	return status
}
