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
	inVerseNum int
	scriptText []string
}

type Stack []string

var hasStyle = map[string]bool{
	`book`: true, `para`: true, `cell`: true, `ms`: true, `note`: true, `sidebar`: true, `figure`: true}

var okUSFM = []string{`p`, // normal paragraph
	`m`,                       // margin paragraph
	`po`,                      // opening of an epistle
	`pr`,                      // right aligned
	`cls`,                     // closure
	`pmo`,                     // embedded text
	`pm`,                      // embedded text paragraph
	`pmc`,                     // embedded text closing
	`pmr`,                     // embedded text refrain
	`pi`, `pi1`, `pi2`, `pi3`, // indented paragraph
	`mi`,                      // indented flush left
	`nb`,                      // no break
	`pb`,                      // explicit page break
	`pc`,                      // centered paragraph
	`ph`, `ph1`, `ph2`, `ph3`, // indented with hanging indent
	`q`, `q1`, `q2`, `q3`, // poetry
	`qr`,                      // right aligned poetry
	`qc`,                      // centered poetry
	`qm`, `qm1`, `qm2`, `qm3`, // embedded text poetic line
	`b`,                       // blank line
	`lh`,                      // list header
	`li`, `li1`, `li2`, `li3`, // list entry
	`lf`,                          // list footer
	`lim`, `lim1`, `lim2`, `lim3`, // embedded list entry
	`mt`, `mt1`, `mt2`, `mt3`, // book title
	`h`, // short book title
}

var numericPattern = regexp.MustCompile(`^\d+`)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:  $HOME/Documents/go2/bin/usx  bibleId")
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
				records := decode(filename)
				records2 := filterByUSFM(records)
				records3 := addChapterHeading(records2)
				loadDatabase(db, records3)
			}
		}
	}
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
	var inVerseNum int
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
				inVerseNum = findIntAttr(se, `number`)
			}
			//if tagName != `char` && tagName != `chapter` && tagName != `verse` {
			if hasStyle[tagName] {
				stack = stack.Push(usfmStyle)
				usfmStyle = findAttr(se, `style`)
			}
			// for version 3.0 add the following
			// find chapter eid to unset chapter number
			// find verse eid to unset verse number
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
			rec.audioFile = tagName
		case xml.EndElement:
			if hasStyle[se.Name.Local] {
				stack, usfmStyle = stack.Pop()
			}
		}
		if bookId != rec.bookId || chapterNum != rec.chapterNum || inVerseNum != rec.inVerseNum ||
			usfmStyle != rec.usfmStyle {
			if rec.bookId != `` && len(rec.scriptText) > 0 {
				records = append(records, rec)
			}
			scriptNum := rec.scriptNum + 1
			if bookId != rec.bookId || chapterNum != rec.chapterNum {
				scriptNum = 1
			}
			rec = ScriptRec{bookId: bookId, chapterNum: chapterNum, scriptNum: scriptNum,
				inVerseNum: inVerseNum, usfmStyle: usfmStyle}
		}
	}
	fmt.Println("num records", len(records))
	return records
}

func filterByUSFM(records []ScriptRec) []ScriptRec {
	var usfmMap = make(map[string]bool)
	for _, usfm := range okUSFM {
		usfmMap[usfm] = true
	}
	var results = make([]ScriptRec, 0, len(records))
	for _, rec := range records {
		_, ok := usfmMap[rec.usfmStyle]
		if ok {
			results = append(results, rec)
		}
	}
	return results
}

func addChapterHeading(records []ScriptRec) []ScriptRec {
	const (
		begin = iota + 1
		newBook
		mtUSFM
	)
	var results = make([]ScriptRec, 0, len(records))
	var shortTitle string
	var lastBookId = ``
	var lastChapter = 0
	var scriptNum = 0
	var state = begin
	for _, rec := range records {
		if rec.usfmStyle == `h` {
			shortTitle = strings.Join(rec.scriptText, ``)
		}
		if state == begin {
			if rec.bookId != lastBookId {
				scriptNum = 1
				state = newBook
			} else if rec.chapterNum != lastChapter {
				var rec2 = rec
				scriptNum = 1
				rec2.scriptNum = scriptNum
				rec2.inVerseNum = 0
				rec2.scriptText = []string{shortTitle + " " + strconv.Itoa(rec.chapterNum)}
				results = append(results, rec2)
			}
			lastBookId = rec.bookId
			lastChapter = rec.chapterNum
		} else if state == newBook {
			if strings.HasPrefix(rec.usfmStyle, `mt`) {
				state = mtUSFM
			} else {
				fmt.Println("Error: after book_id change, mt expected, but,", rec)
			}
		} else if state == mtUSFM {
			if !strings.HasPrefix(rec.usfmStyle, `mt`) {
				var rec2 = rec
				scriptNum += 1
				rec2.scriptNum = scriptNum
				rec2.inVerseNum = 0
				rec2.scriptText = []string{shortTitle + " " + strconv.Itoa(rec.chapterNum)}
				results = append(results, rec2)
				state = begin
			}
		}
		if rec.usfmStyle != `h` {
			scriptNum += 1
			rec.scriptNum = scriptNum
			results = append(results, rec)
		}
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
			in_verse_num INTEGER,
			script_text TEXT,
			script_begin_ts REAL,
			script_end_ts REAL,
			script_mfcc BLOB,
			mfcc_rows INTEGER,
			mfcc_cols INTEGER) STRICT`
	_, err := db.Exec(sqlStmt)
	//fmt.Println("result", result)
	if err != nil {
		log.Fatal("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt = `INSERT INTO audio_scripts(book_id, chapter_num, audio_file, 
			script_num, usfm_style, person, actor, in_verse_num, script_text) 
			VALUES (?,?,?,?,?,?,?,?,?)`
	stmt, err := tx.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, rec := range records {
		text := strings.Join(rec.scriptText, ``)
		_, err = stmt.Exec(rec.bookId, rec.chapterNum, rec.audioFile, rec.scriptNum,
			rec.usfmStyle, rec.person, rec.actor, rec.inVerseNum, text)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

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
		fmt.Println(count)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
