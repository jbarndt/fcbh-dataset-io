package read

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"github.com/xuri/excelize/v2"
	"strconv"
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

func (r ScriptReader) ProcessFiles(files []input.InputFile) dataset.Status {
	var status dataset.Status
	for _, file := range files {
		status = r.Read(file.FilePath())
	}
	return status
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
