package read

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// This program will read Excel data and load the audio_scripts table

type ScriptReader struct {
	ctx context.Context
	db  db.DBAdapter
}

func NewScriptReader(db db.DBAdapter) ScriptReader {
	var d ScriptReader
	d.ctx = db.Ctx
	d.db = db
	return d
}

func (r ScriptReader) FindFile(bibleId string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	directory := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	files, err := os.ReadDir(directory)
	if err != nil {
		status = log.Error(r.ctx, 500, err, "Could not read directory", directory)
		return result, status
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".xlsx") {
			return filepath.Join(directory, filename), status
		}
	}
	status = log.Error(r.ctx, 500, err, "Could not find .xlsx file in", directory)
	return result, status
}

func (r ScriptReader) Read(filePath string) dataset.Status {
	var status dataset.Status
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return log.Error(r.ctx, 500, err, "Error: could not open", filePath)
	}
	defer file.Close()
	sheets := file.GetSheetList()
	sheet := sheets[0]
	rows, err := file.GetRows(sheet)
	if err != nil {
		return log.Error(r.ctx, 500, err, `Error reading excel file.`)
	}
	var records []db.Script
	for i, row := range rows {
		if i == 0 {
			continue // skip headings
		}
		var rec db.Script
		switch row[1] {
		case `JMS`:
			rec.BookId = `JAS`
		case `TTS`:
			rec.BookId = `TIT`
		case ``:
			return log.ErrorNoErr(r.ctx, 500, `Error: Did not find book_id`)
		default:
			rec.BookId = row[1]
		}
		rec.ChapterNum, err = strconv.Atoi(row[2])
		if err != nil {
			return log.Error(r.ctx, 500, err, "Error: chapter num is not numeric", row[2])
		}
		if row[3] == `<<` {
			rec.VerseStr = ``
			rec.VerseNum = 0
		} else {
			rec.VerseStr = row[3]
			rec.VerseNum, err = strconv.Atoi(row[3])
			if err != nil {
				return log.Error(r.ctx, 500, err, `Error: verse num is not numeric`, row[3])
			}
		}
		rec.Person = row[4]
		//rec.Actor = row[5]
		rec.ScriptNum = row[5]
		text := row[8]
		//text = strings.Replace(text,'_x000D_','' ) // remove excel CR
		rec.ScriptTexts = append(rec.ScriptTexts, text)
		if rec.ScriptNum[len(rec.ScriptNum)-1] != 'r' {
			records = append(records, rec)
		}
	}
	status = r.db.InsertScripts(records)
	return status
}
