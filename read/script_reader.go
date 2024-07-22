package read

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"github.com/xuri/excelize/v2"
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
	var col ColIndex
	var records []db.Script
	for i, row := range rows {
		if i == 0 {
			col, status = r.FindColIndexes(row)
			if status.IsErr {
				return status
			}
			continue
		}
		var rec db.Script
		switch row[col.BookCol] {
		case `JMS`:
			rec.BookId = `JAS`
		case `TTS`:
			rec.BookId = `TIT`
		case ``:
			return log.ErrorNoErr(r.ctx, 500, `Error: Did not find book_id`)
		default:
			rec.BookId = row[col.BookCol]
		}
		rec.ChapterNum, err = strconv.Atoi(row[col.ChapterCol])
		if err != nil {
			return log.Error(r.ctx, 500, err, "Error: chapter num is not numeric", row[col.ChapterCol])
		}
		if row[col.VerseCol] == `<<` {
			rec.VerseStr = `0`
			rec.VerseNum = 0
		} else {
			rec.VerseStr = row[col.VerseCol]
			rec.VerseNum, err = strconv.Atoi(row[col.VerseCol])
			if err != nil {
				return log.Error(r.ctx, 500, err, `Error: verse num is not numeric`, row[3])
			}
		}
		rec.Person = row[col.CharacterCol]
		rec.ScriptNum = row[col.LineCol]
		text := row[col.TextCol]
		//text = strings.Replace(text,'_x000D_','' ) // remove excel CR
		rec.ScriptTexts = append(rec.ScriptTexts, text)
		if rec.ScriptNum[len(rec.ScriptNum)-1] != 'r' {
			records = append(records, rec)
		}
	}
	status = r.db.InsertScripts(records)
	return status
}

type ColIndex struct {
	BookCol      int
	ChapterCol   int
	VerseCol     int
	CharacterCol int
	LineCol      int
	TextCol      int
}

func (r ScriptReader) FindColIndexes(heading []string) (ColIndex, dataset.Status) {
	var c ColIndex
	for col, head := range heading {
		switch strings.ToLower(head) {
		case `book`, `bk`:
			c.BookCol = col
		case `chapter`, `cp`:
			c.ChapterCol = col
		case `verse`, `verse_number`:
			c.VerseCol = col
		case `line_number`, `line id`, `line`:
			c.LineCol = col
		case `characters1`, `character`:
			c.CharacterCol = col
		case `verse_content1`, `target language`:
			c.TextCol = col
		}
	}
	var msgs []string
	if c.BookCol == 0 {
		msgs = append(msgs, `Book column was not found`)
	}
	if c.ChapterCol == 0 {
		msgs = append(msgs, `Chapter column was not found`)
	}
	if c.VerseCol == 0 {
		msgs = append(msgs, `Verse column was not found`)
	}
	if c.LineCol == 0 {
		msgs = append(msgs, `Line column was not found`)
	}
	if c.TextCol == 0 {
		msgs = append(msgs, `Text column was not found`)
	}
	var status dataset.Status
	if len(msgs) > 0 {
		status = log.ErrorNoErr(r.ctx, 500, `Columns missing in script`, strings.Join(msgs, `; `))
	}
	return c, status
}
