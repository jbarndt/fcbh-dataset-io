package read

import (
	"dataset_io"
	"dataset_io/db"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DBPTextReader struct {
	conn db.DBAdapter
}

func NewDBPTextReader(conn db.DBAdapter) *DBPTextReader {
	return &DBPTextReader{conn: conn}
}

func (d *DBPTextReader) ProcessDirectory(bibleId string, testament dataset_io.TestamentType) {
	directory := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	switch testament {
	case dataset_io.NT:
		d.processFile(directory, bibleId+"N_ET.json")
	case dataset_io.OT:
		d.processFile(directory, bibleId+"O_ET.json")
	case dataset_io.ONT:
		d.processFile(directory, bibleId+"O_ET.json")
		d.processFile(directory, bibleId+"N_ET.json")
	default:
		log.Fatal("Error: unknown TestamentType", testament)
	}
}

func (d *DBPTextReader) processFile(directory, filename string) {
	var scriptNum = 0
	var lastBookId string
	filePath := filepath.Join(directory, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", filePath, err)
		return
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
		fmt.Println("Error parsing JSON:", err)
		return
	}
	fmt.Println("num verses", len(verses))
	var records = make([]db.InsertScriptRec, 0, 1000)
	for _, vs := range verses {
		scriptNum++
		if vs.BookId != lastBookId {
			fmt.Println(vs.BookId)
			lastBookId = vs.BookId
			scriptNum = 1
		}
		var rec db.InsertScriptRec
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
		rec.ScriptText = append(rec.ScriptText, text)
		records = append(records, rec)
	}
	d.conn.InsertScripts(records)
}
