package read

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DBPTextReader struct {
	ctx  context.Context
	conn db.DBAdapter
}

func NewDBPTextReader(conn db.DBAdapter) DBPTextReader {
	var d DBPTextReader
	d.ctx = conn.Ctx
	d.conn = conn
	return d
}

func (d *DBPTextReader) ProcessDirectory(bibleId string, testament dataset.TestamentType) dataset.Status {
	var status dataset.Status
	directory := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	switch testament {
	case dataset.NT:
		d.processFile(directory, bibleId+"N_ET.json")
	case dataset.OT:
		d.processFile(directory, bibleId+"O_ET.json")
	case dataset.C:
		d.processFile(directory, bibleId+"O_ET.json")
		d.processFile(directory, bibleId+"N_ET.json")
	default:
		status = log.ErrorNoErr(d.ctx, 500, "Error: unknown TestamentType", testament)
	}
	return status
}

func (d *DBPTextReader) processFile(directory, filename string) dataset.Status {
	var status dataset.Status
	var scriptNum = 0
	var lastBookId string
	filePath := filepath.Join(directory, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return log.Error(d.ctx, 500, err, "Error reading file:", filePath)
	}
	fmt.Println("Read", filename, len(content), "bytes")
	type TempRec struct {
		BookId     string `json:"book_id"`
		ChapterNum int    `json:"chapter"`
		VerseStart int    `json:"verse_start"`
		VerseEnd   int    `json:"verse_end"`
		Text       string `json:"verse_text"`
	}
	var verses []TempRec
	err = json.Unmarshal(content, &verses)
	if err != nil {
		return log.Error(d.ctx, 500, err, "Error parsing JSON from plain_text")
	}
	fmt.Println("num verses", len(verses))
	var records = make([]db.Script, 0, 1000)
	for _, vs := range verses {
		scriptNum++
		if vs.BookId != lastBookId {
			fmt.Println(vs.BookId)
			lastBookId = vs.BookId
			scriptNum = 1
		}
		var rec db.Script
		rec.ScriptNum = strconv.Itoa(scriptNum)
		rec.BookId = vs.BookId
		rec.ChapterNum = vs.ChapterNum
		rec.VerseNum = vs.VerseStart
		if vs.VerseStart == vs.VerseEnd {
			rec.VerseStr = strconv.Itoa(vs.VerseStart)
		} else {
			rec.VerseStr = strconv.Itoa(vs.VerseStart) + `-` + strconv.Itoa(vs.VerseEnd)
		}
		text := strings.Replace(vs.Text, "&lt", "<", -1)
		text = strings.Replace(text, "&gt", ">", -1)
		rec.ScriptTexts = append(rec.ScriptTexts, text)
		records = append(records, rec)
	}
	status = d.conn.InsertScripts(records)
	return status
}
