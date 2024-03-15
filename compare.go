package main

import (
	"database/sql"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	database1, database2 := getCommandLine()
	fmt.Println("1", database1, "2", database2)
	var db1 = openDatabase(database1)
	var db2 = openDatabase(database2)
	var numChapters = readNumChapters(db1)
	var nt = []string{`MAT`, `MRK`, `LUK`, `JHN`, `ACT`, `ROM`, `1CO`, `2CO`, `GAL`, `EPH`, `PHP`, `COL`,
		`1TH`, `2TH`, `1TI`, `2TI`, `TIT`, `PHM`, `HEB`, `JAS`, `1PE`, `2PE`, `1JN`, `2JN`, `3JN`, `JUD`, `REV`}
	for _, bookId := range nt {
		var chapInBook, _ = numChapters[bookId]
		var chapter = 1
		for chapter <= chapInBook {
			lines1 := process(db1, database1, bookId, chapter)
			lines2 := process(db2, database2, bookId, chapter)
			diff(lines1, lines2)
			chapter++
		}
	}
	db1.Close()
	db2.Close()
	//displayTest(lines1)
	//displayTest(lines2)
}

func getCommandLine() (string, string) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: $HOME/Documents/go2/bin/compare  database1  database2")
		os.Exit(1)
	}
	return os.Args[1], os.Args[2]
}

func openDatabase(name string) *sql.DB {
	path := os.Getenv(`FCBH_DATASET_DB`)
	database := filepath.Join(path, name)
	db, err := sql.Open("sqlite", database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func readNumChapters(db *sql.DB) map[string]int {
	var results = make(map[string]int)
	sqlStmt := `SELECT book_id, max(chapter_num) FROM audio_scripts 
			GROUP BY book_id`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	type rec struct {
		bookId      string
		numChapters int
	}
	for rows.Next() {
		var tmp rec
		err := rows.Scan(&tmp.bookId, &tmp.numChapters)
		if err != nil {
			log.Fatal(err)
		}
		results[tmp.bookId] = tmp.numChapters
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return results
}

func process(db *sql.DB, database string, bookId string, chapterNum int) []Verse {
	lines := readDatabase(db, bookId, chapterNum)
	//lines = groupBySentence(lines)
	if strings.Contains(database, "SCRIPT") {
		lines = consolidateScript(lines)
	} else if strings.Contains(database, "USX") {
		lines = consolidateUSX(lines)
	}
	return lines
}

type Verse struct {
	bookId  string
	chapter int
	num     string
	text    string
}

func readDatabase(db *sql.DB, bookId string, chapterNum int) []Verse {
	type tmpVerse struct {
		bookId  string
		chapter int
		num     sql.NullInt32
		text    sql.NullString
	}
	sqlStmt := `SELECT book_id, chapter_num, in_verse_num, script_text FROM audio_scripts 
			WHERE book_id=? AND chapter_num=?
			ORDER BY script_id`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(bookId, chapterNum)
	if err != nil {
		log.Fatal(err)
	}
	var results []Verse
	for rows.Next() {
		var tmp tmpVerse
		err := rows.Scan(&tmp.bookId, &tmp.chapter, &tmp.num, &tmp.text)
		if err != nil {
			log.Fatal(err)
		}
		if tmp.text.Valid {
			var verse = Verse{bookId: tmp.bookId, chapter: tmp.chapter,
				num: strconv.Itoa(int(tmp.num.Int32)), text: tmp.text.String}
			results = append(results, verse)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return results
}

func consolidateScript(verses []Verse) []Verse {
	const (
		begin = iota + 1
		inNum
		endNum
	)
	//var labels = []string{``, `BEGIN`, `INNUM`, `ENDNUM`}
	var results = make([]Verse, 0, len(verses))
	var sumInput = 0
	var sumOutput = 0
	// Somehow do this one chapter at a time.
	var bookId = verses[0].bookId
	var chapter = verses[0].chapter
	var parts = make([]string, 0, 100)
	for _, rec := range verses {
		parts = append(parts, rec.text)
		sumInput += len(rec.text)
	}
	text := strings.Join(parts, ``)
	//fmt.Println(text, "\n")
	var verseNum = `0`
	var tmpNum []byte
	var index = 0
	var state = begin
	for index < len(text) {
		//fmt.Println(`state:`, labels[state], `index:`, index, `char:`, string(text[index]))
		switch state {
		case begin:
			var part string
			search := text[index:]
			pos := strings.Index(search, `{`)
			if pos < 0 {
				part = search
				sumOutput += len(part)
			} else {
				part = search[:pos]
				state = inNum
				tmpNum = []byte{}
				sumOutput += len(part) + 1
			}
			verse := Verse{bookId: bookId, chapter: chapter, num: verseNum, text: part}
			results = append(results, verse)
			//fmt.Println(part)
			index += len(part) + 1
		case inNum:
			char := text[index]
			//fmt.Println("inNum", string(char))
			if char >= '0' && char <= '9' {
				tmpNum = append(tmpNum, char)
			} else if char == '}' {
				verseNum = string(tmpNum)
				state = endNum
			} else {
				logError(bookId, chapter, verseNum, `Invalid char in {nn, expect n or }`,
					string(char))
			}
			index++
			sumOutput += 1
		case endNum:
			char := text[index]
			peek := text[index+1]
			//fmt.Println("endNum", string(char), string(peek))
			if (char == '_' || char == '-') && peek == '{' {
				state = inNum
				tmpNum = []byte(verseNum + "-")
				index += 2
				sumOutput += 2
			} else {
				state = begin
			}
		}
	}
	if sumInput != sumOutput {
		log.Fatal("Bug: Not all data processed by consolidateScript input:", sumInput, " output:", sumOutput)
	}
	return results
}

func logError(bookId string, chapter int, verseNum string, message string, found string) {
	fmt.Println("Error:", bookId, chapter, verseNum, message, found)
	os.Exit(1)
}

func consolidateUSX(verses []Verse) []Verse {
	var sumInput = 0
	var sumOutput = 0
	var results = make([]Verse, 0, len(verses))
	var lastChapter = -1
	var lastVerse = ``
	var verse Verse
	for _, rec := range verses {
		sumInput += len(rec.text)
		if rec.chapter != lastChapter || rec.num != lastVerse {
			if lastChapter != -1 {
				results = append(results, verse)
				sumOutput += len(verse.text)
			}
			lastChapter = rec.chapter
			lastVerse = rec.num
			verse = Verse{bookId: rec.bookId, chapter: rec.chapter, num: rec.num, text: ``}
		}
		verse.text += rec.text
	}
	if len(verse.text) > 0 {
		results = append(results, verse)
		sumOutput += len(verse.text)
	}
	if sumInput != sumOutput {
		log.Fatal("Bug: Not all data processed by consolidateUSX input:", sumInput, " output:", sumOutput)
	}
	return results
}

func groupBySentence(verses []Verse) []Verse {
	var results = make([]Verse, 0, len(verses))
	var text []string
	for _, verse := range verses {
		text = append(text, verse.text)
	}
	chapter := strings.Join(text, ``)
	parts := strings.Split(chapter, `.`)
	for i, part := range parts {
		verse := Verse{bookId: `MAT`, chapter: 2, num: strconv.Itoa(i), text: part + `.`}
		results = append(results, verse)
	}
	return results
}

func displayTest(verses []Verse) {
	for i, rec := range verses {
		fmt.Println(i, rec.bookId, rec.chapter, rec.num, len(rec.text))
	}
	fmt.Println("========")
}

/* This diff method assumes one chapter at a time */
func diff(verses1 []Verse, verses2 []Verse) {
	type pair struct {
		bookId  string
		chapter int
		num     string
		text1   string
		text2   string
	}
	// Put the second data in a map
	var verse2Map = make(map[string]Verse)
	for _, vs2 := range verses2 {
		verse2Map[vs2.num] = vs2
	}
	// combine the verse2 to verse1 that match
	var didMatch = make(map[string]bool)
	var pairs = make([]pair, 0, len(verses1))
	for _, vs1 := range verses1 {
		vs2, ok := verse2Map[vs1.num]
		if ok {
			didMatch[vs1.num] = true
		}
		p := pair{bookId: vs1.bookId, chapter: vs1.chapter, num: vs1.num, text1: vs1.text, text2: vs2.text}
		pairs = append(pairs, p)
	}
	// pick up any verse2 that did not match verse1
	for _, vs2 := range verses2 {
		_, ok := didMatch[vs2.num]
		if !ok {
			p := pair{bookId: vs2.bookId, chapter: vs2.chapter, num: vs2.num, text1: ``, text2: vs2.text}
			pairs = append(pairs, p)
		}
	}
	// perform a match on pairs
	diffMatch := diffmatchpatch.New()
	for _, pair := range pairs {
		if len(pair.text1) > 0 || len(pair.text2) > 0 {
			diffs := diffMatch.DiffMain(pair.text1, pair.text2, false)
			if !isMatch(diffs) {
				ref := pair.bookId + " " + strconv.Itoa(pair.chapter) + ":" + pair.num
				fmt.Println(ref, diffMatch.DiffPrettyText(diffs))
				fmt.Println("=============")
			}
		}
	}
}

func isMatch(diffs []diffmatchpatch.Diff) bool {
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert || diff.Type == diffmatchpatch.DiffDelete {
			if len(strings.TrimSpace(diff.Text)) > 0 {
				return false
			}
		}
	}
	return true
}
