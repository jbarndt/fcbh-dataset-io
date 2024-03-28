package read

import (
	"dataset/db"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// This program will read Excel data and load the audio_scripts table

type ScriptReader struct {
	db db.DBAdapter
}

func NewScriptReader(db db.DBAdapter) ScriptReader {
	var d ScriptReader
	d.db = db
	return d
}

func (r ScriptReader) FindFile(bibleId string) string {
	directory := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal("Could not read directory", err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".xlsx") {
			return filepath.Join(directory, filename)
		}
	}
	log.Fatalln("Could not find .xlsx file in", directory)
	return ``
}

func (r ScriptReader) Read(filePath string) {
	fmt.Println("reading", filePath)
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("Error: could not open", filePath, err)
		return
	}
	defer file.Close()
	sheets := file.GetSheetList()
	sheet := sheets[0]
	rows, err := file.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}
	var records []db.InsertScriptRec
	for i, row := range rows {
		if i == 0 {
			continue // skip headings
		}
		var rec db.InsertScriptRec
		switch row[1] {
		case `JMS`:
			rec.BookId = `JAS`
		case `TTS`:
			rec.BookId = `TIT`
		case ``:
			log.Fatalln(`Error: Did not find book_id`)
		default:
			rec.BookId = row[1]
		}
		rec.ChapterNum, err = strconv.Atoi(row[2])
		if err != nil {
			log.Fatalln("Error: chapter num is not numeric", row[2])
		}
		if row[3] == `<<` {
			rec.VerseStr = ``
			rec.VerseNum = 0
		} else {
			rec.VerseStr = row[3]
			rec.VerseNum, err = strconv.Atoi(row[3])
			if err != nil {
				log.Fatalln(`Error: verse num is not numeric`, row[3])
			}
		}
		rec.Person = row[4]
		//rec.Actor = row[5]
		rec.ScriptNum = row[5]
		text := row[8]
		//text = strings.Replace(text,'_x000D_','' ) // remove excel CR
		rec.ScriptText = append(rec.ScriptText, text)
		if rec.ScriptNum[len(rec.ScriptNum)-1] != 'r' {
			records = append(records, rec)
		}
	}
	r.db.InsertScripts(records)
}
