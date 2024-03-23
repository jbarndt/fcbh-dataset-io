// package usx
package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ScriptRec struct {
	bookId     string
	chapterNum int
	audioFile  string
	scriptNum  int
	usfmStyle  string
	person     string
	actor      string
	verseNum   int
	verseStr   string
	scriptText []string
}

type Stack []string

var hasStyle = map[string]bool{
	`book`: true, `para`: true, `char`: true, `cell`: true, `ms`: true, `note`: true, `sidebar`: true, `figure`: true}

var numericPattern = regexp.MustCompile(`^\d+`)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:  $HOME/Documents/go2/bin/usx_parser  bibleId")
		os.Exit(1)
	}
	var bibleId = os.Args[1]
	dbPath := os.Getenv(`FCBH_DATASET_DB`)
	var db = openDatabase(dbPath, bibleId+"_USXEDIT.db")
	directory := filepath.Join(dbPath, `download`, bibleId)
	dirs, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		if strings.HasSuffix(dir.Name(), `N_ET-usx`) { // New Testament only
			//if strings.HasSuffix(dir.Name(), `_ET-usx`) {
			subDir := filepath.Join(directory, dir.Name())
			files, err := os.ReadDir(subDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				filename := filepath.Join(subDir, file.Name())
				fmt.Println(filename)
				records := decode(filename) // Also edits out non-script elements
				titleDesc := extractTitles(records)
				records = addChapterHeading(records, titleDesc)
				records = correctScriptNum(records)
				loadDatabase(db, records)
			}
		}
	}
	countRecords(db)
	db.Close()
}

func decode(filename string) []ScriptRec {
	xmlFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	var stack Stack
	var rec ScriptRec
	var records []ScriptRec
	var tagName string
	var bookId string
	var chapterNum = 1
	var verseNum int
	var verseStr string
	var usfmStyle string
	decoder := xml.NewDecoder(xmlFile)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			panic(err)
		}
		switch se := token.(type) {
		//[StartElement], [EndElement], [CharData], [Comment], [ProcInst], [Directive].
		case xml.StartElement:
			tagName = se.Name.Local
			if tagName == `book` {
				bookId = findAttr(se, `code`)
			} else if tagName == `chapter` {
				chapterNum = findIntAttr(se, `number`)
			} else if tagName == `verse` {
				verseNum = findIntAttr(se, `number`)
				verseStr = findAttr(se, `number`)
			}
			if hasStyle[tagName] {
				usfmStyle = tagName + `.` + findAttr(se, `style`)
				if include(usfmStyle) { // This if fileters out the non-script code
					stack = stack.Push(usfmStyle)
				} else {
					err := decoder.Skip()
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case xml.CharData:
			text := string(se)
			if len(strings.TrimSpace(text)) > 0 {
				// This is needed because scripts use {n} as a verse number.
				if strings.Contains(text, "{") || strings.Contains(text, "}") {
					fmt.Println("Found { } changed to ( ) in:", text)
					text = strings.Replace(text, `{`, `(`, -1)
					text = strings.Replace(text, `}`, `)`, -1)
				}
				rec.scriptText = append(rec.scriptText, text)
			}
		case xml.EndElement:
			if hasStyle[se.Name.Local] {
				stack, usfmStyle = stack.Pop()
			}
		}
		if chapterNum != rec.chapterNum || verseNum != rec.verseNum || usfmStyle != rec.usfmStyle {
			if bookId != `` && len(rec.scriptText) > 0 {
				records = append(records, rec)
			}
			scriptNum := rec.scriptNum + 1
			if chapterNum != rec.chapterNum {
				scriptNum = 1
			}
			rec = ScriptRec{bookId: bookId, chapterNum: chapterNum, scriptNum: scriptNum,
				verseNum: verseNum, verseStr: verseStr, usfmStyle: usfmStyle}
		}
	}
	fmt.Println("num records", len(records))
	return records
}

type titleDesc struct {
	heading string
	title   []ScriptRec
	areDiff bool
}

func extractTitles(records []ScriptRec) titleDesc {
	var results titleDesc
	for _, rec := range records {
		if rec.usfmStyle == `para.h` {
			results.heading = strings.Join(rec.scriptText, ``)
		} else if strings.HasPrefix(rec.usfmStyle, `para.mt`) {
			results.title = append(results.title, rec)
		}
	}
	if results.heading == `` && len(results.title) > 0 {
		results.heading = strings.Join(results.title[len(results.title)-1].scriptText, ``)
	}
	return results
}

func addChapterHeading(records []ScriptRec, titles titleDesc) []ScriptRec {
	var results = make([]ScriptRec, 0, len(records))
	for _, rec := range titles.title {
		results = append(results, rec)
	}
	var lastChapter = -1
	for _, rec := range records {
		if rec.chapterNum != lastChapter {
			lastChapter = rec.chapterNum
			var rec2 = rec
			rec2.verseNum = 0
			rec2.verseStr = ``
			rec2.scriptText = []string{titles.heading + " " + strconv.Itoa(rec.chapterNum)}
			results = append(results, rec2)
		}
		if rec.usfmStyle != `para.h` && !strings.HasPrefix(rec.usfmStyle, `para.mt`) {
			results = append(results, rec)
		}
	}
	return results
}

func correctScriptNum(records []ScriptRec) []ScriptRec {
	var results = make([]ScriptRec, 0, len(records))
	var scriptNum = 0
	var lastChapter = 0
	for _, rec := range records {
		if rec.chapterNum != lastChapter {
			lastChapter = rec.chapterNum
			scriptNum = 0
		}
		scriptNum += 1
		rec.scriptNum = scriptNum
		results = append(results, rec)
	}
	return results
}

func findAttr(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ``
}

func findIntAttr(se xml.StartElement, name string) int {
	val := findAttr(se, name)
	if val == `` {
		return 0
	} else {
		num := numericPattern.FindString(val)
		result, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (s Stack) Push(v string) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, string) {
	l := len(s)
	if l < 1 {
		panic("You tried to pop an empty stack.")
	}
	return s[:l-1], s[l-1]
}

func openDatabase(path, name string) *sql.DB {
	database := filepath.Join(path, name)
	os.Remove(database)
	db, err := sql.Open("sqlite", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func loadDatabase(db *sql.DB, records []ScriptRec) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS audio_scripts (
			script_id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id TEXT NOT NULL,
			chapter_num INTEGER NOT NULL,
			audio_file TEXT NOT NULL,
			script_num TEXT NOT NULL,
			usfm_style TEXT,
			person TEXT,  
			actor TEXT,  
			verse_num INTEGER NOT NULL,
			verse_str TEXT NOT NULL,
			script_text TEXT NOT NULL,
			script_begin_ts REAL,
			script_end_ts REAL,
			script_mfcc BLOB,
			mfcc_rows INTEGER,
			mfcc_cols INTEGER) STRICT`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt = `INSERT INTO audio_scripts(book_id, chapter_num, audio_file, 
			script_num, usfm_style, person, actor, verse_num, verse_str, script_text) 
			VALUES (?,?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, rec := range records {
		text := strings.Join(rec.scriptText, ``)
		_, err = stmt.Exec(rec.bookId, rec.chapterNum, rec.audioFile, rec.scriptNum,
			rec.usfmStyle, rec.person, rec.actor, rec.verseNum, rec.verseStr, text)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func countRecords(db *sql.DB) {
	rows, err := db.Query(`select count(*) from audio_scripts`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Total Records", count)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
