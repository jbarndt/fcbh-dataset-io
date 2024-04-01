package db

import (
	"context"
	"database/sql"
	log "dataset/logger"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
	"time"
)

func TestGenericQuery(t *testing.T) {
	ctx := context.Background()
	var database = `ATIWBT_USXEDIT.db`
	conn := NewDBAdapter(ctx, database)
	query := SelectType{conn.DB}
	sql1 := `SELECT book_id, chapter_num, script_begin_ts FROM scripts WHERE chapter_num=?`
	results, err := query.Select(sql1, 11)
	if err != nil {
		log.Fatal(ctx, err)
	}
	for _, result := range results {
		fmt.Println(result[0].(string), result[1].(int), result[2].(float64))
	}
}

var test1 = `SELECT i.dataset_id, i.bible_id, i.audio_fileset_id, i.text_fileset_id, i.text_source, i.language_iso,
i.version_code, i.languge_id, s.script_id, s.book_id, s.chapter_num, s.audio_file,
s.script_num, s.usfm_style, s.person, s.actor, s.verse_str, s.script_text,
w.word_id, w.word_seq, w.verse_num, w.ttype, w.word
FROM ident i 
JOIN scripts s ON i.dataset_id = s.dataset_id
JOIN words w ON s.script_id = w.script_id
ORDER BY w.word_id`

var test2 = `SELECT s.script_id, s.book_id, s.chapter_num, s.audio_file,
s.script_num, s.usfm_style, s.person, s.actor, s.verse_str, s.script_text,
w.word_id, w.word_seq, w.verse_num, w.ttype, w.word
FROM scripts s 
JOIN words w ON s.script_id = w.script_id
ORDER BY w.word_id`

type Test1Rec struct {
	DatasetId      int
	BibleId        string
	AudioFilesetId string
	TextFilesetId  string
	TextSource     string
	LanguageIso    string
	VersionCode    string
	LanguageId     int
	ScriptId       int
	BookId         string
	ChapterNum     int
	AudioFile      string
	ScriptNum      string
	UsfmStyle      string
	Person         string
	Actor          string
	VerseStr       string
	ScriptText     string
	WordId         int
	WordSeq        int
	VerseNum       int
	Ttype          string
	Word           string
}

func TestStandardInterface(t *testing.T) {
	var start = time.Now()
	ctx := context.Background()
	var db = NewDBAdapter(ctx, `BGGWFW_SCRIPT.db`)
	rows, err := db.DB.Query(test2)
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	defer rows.Close()
	var result = make([]Test1Rec, 0, 500000)
	for rows.Next() {
		var rec Test1Rec
		err := rows.Scan(
			/*
				&rec.DatasetId,
				&rec.BibleId,
				&rec.AudioFilesetId,
				&rec.TextFilesetId,
				&rec.TextSource,
				&rec.LanguageIso,
				&rec.VersionCode,
				&rec.LanguageId,*/
			&rec.ScriptId,
			&rec.BookId,
			&rec.ChapterNum,
			&rec.AudioFile,
			&rec.ScriptNum,
			&rec.UsfmStyle,
			&rec.Person,
			&rec.Actor,
			&rec.VerseStr,
			&rec.ScriptText,
			&rec.WordId,
			&rec.WordSeq,
			&rec.VerseNum,
			&rec.Ttype,
			&rec.Word)
		if err != nil {
			log.Fatal(ctx, err, test1)
		}
		result = append(result, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	fmt.Println("Count", len(result))
	fmt.Println("Elapsed", time.Since(start))
}

func TestGenericInterface(t *testing.T) {
	var start = time.Now()
	ctx := context.Background()
	var db = NewDBAdapter(ctx, `BGGWFW_SCRIPT.db`)
	var query = SelectType{db.DB}
	records, err := query.Select(test1)
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	fmt.Println("Query Elapsed", time.Since(start))
	var results = make([]Test1Rec, 0, 500000)
	for _, s := range records {
		var rec Test1Rec
		rec.DatasetId = s[0].(int)
		rec.BibleId = s[1].(string)
		rec.AudioFilesetId = s[2].(string)
		rec.TextFilesetId = s[3].(string)
		rec.TextSource = s[4].(string)
		rec.LanguageIso = s[5].(string)
		rec.VersionCode = s[6].(string)
		rec.LanguageId = s[7].(int)
		rec.ScriptId = s[8].(int)
		rec.BookId = s[9].(string)
		rec.ChapterNum = s[10].(int)
		rec.AudioFile = s[11].(string)
		rec.ScriptNum = s[12].(string)
		rec.UsfmStyle = s[13].(string)
		rec.Person = s[14].(string)
		rec.Actor = s[15].(string)
		rec.VerseStr = s[16].(string)
		rec.ScriptText = s[17].(string)
		rec.WordId = s[18].(int)
		rec.WordSeq = s[19].(int)
		rec.VerseNum = s[20].(int)
		rec.Ttype = s[21].(string)
		rec.Word = s[22].(string)
		results = append(results, rec)
	}
	fmt.Println("Count", len(results))
	fmt.Println("Total Elapsed", time.Since(start))
}

type ScriptRec struct {
	ScriptId   int
	BookId     string
	ChapterNum int
	AudioFile  string
	ScriptNum  string
	UsfmStyle  string
	Person     string
	Actor      string
	VerseStr   string
	ScriptText string
}
type WordRec struct {
	WordId   int
	WordSeq  int
	VerseNum int
	Ttype    string
	Word     string
}

func TestManyQuery(t *testing.T) {
	//var db = NewDBAdapter(`BGGWFW_SCRIPT.db`)
	ctx := context.Background()
	database := os.Getenv(`FCBH_DATASET_DB`) + "/BGGWFW_SCRIPT.db"
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(ctx, err)
	}
	scripts := scriptsQuery(db)
	fmt.Println(len(scripts))
	var count = 0
	for _, script := range scripts {
		words := wordsQuery(db, script.ScriptId)
		count += len(words)
		//fmt.Println("\t", len(words))
	}
	fmt.Println(count, "Words")
}

func scriptsQuery(db *sql.DB) []ScriptRec {
	var test1 = `SELECT s.script_id, s.book_id, s.chapter_num, s.audio_file,
s.script_num, s.usfm_style, s.person, s.actor, s.verse_str, s.script_text
FROM scripts s ORDER BY s.script_id`
	rows, err := db.Query(test1)
	ctx := context.Background()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	defer rows.Close()
	var results = make([]ScriptRec, 0, 10000)
	for rows.Next() {
		var rec ScriptRec
		err := rows.Scan(
			&rec.ScriptId,
			&rec.BookId,
			&rec.ChapterNum,
			&rec.AudioFile,
			&rec.ScriptNum,
			&rec.UsfmStyle,
			&rec.Person,
			&rec.Actor,
			&rec.VerseStr,
			&rec.ScriptText)
		if err != nil {
			log.Fatal(ctx, err, test1)
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	return results
}

func wordsQuery(db *sql.DB, scriptId int) []WordRec {
	ctx := context.Background()
	var test1 = `SELECT w.word_id, w.word_seq, w.verse_num, w.ttype, w.word
FROM words w WHERE w.script_id = ? ORDER BY w.word_id`
	rows, err := db.Query(test1, scriptId)
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	defer rows.Close()
	var results = make([]WordRec, 0, 100)
	for rows.Next() {
		var rec WordRec
		err := rows.Scan(
			&rec.WordId,
			&rec.WordSeq,
			&rec.VerseNum,
			&rec.Ttype,
			&rec.Word)
		if err != nil {
			log.Fatal(ctx, err, test1)
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	return results
}

type ScriptSmallRec struct {
	ScriptId   int
	BookId     string
	ChapterNum int
	//AudioFile  string
	//ScriptNum  string
	//UsfmStyle  string
	//Person     string
	//Actor      string
	VerseStr string
	//ScriptText string
}
type WordSmallRec struct {
	WordId int
	//WordSeq  int
	//VerseNum int
	//Ttype    string
	Word string
}

func TestManyAmallQuery(t *testing.T) {
	ctx := context.Background()
	//var db = NewDBAdapter(`BGGWFW_SCRIPT.db`)
	database := os.Getenv(`FCBH_DATASET_DB`) + "/BGGWFW_SCRIPT.db"
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(ctx, err)
	}
	scripts := smallScriptsQuery(db)
	fmt.Println(len(scripts))
	var count = 0
	for _, script := range scripts {
		words := smallWordsQuery(db, script.ScriptId)
		count += len(words)
		//fmt.Println("\t", len(words))
	}
	fmt.Println(count, "Words")
}

func smallScriptsQuery(db *sql.DB) []ScriptSmallRec {
	ctx := context.Background()
	var test1 = `SELECT s.script_id, s.book_id, s.chapter_num, s.verse_str
FROM scripts s ORDER BY s.script_id`
	rows, err := db.Query(test1)
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	defer rows.Close()
	var results = make([]ScriptSmallRec, 0, 10000)
	for rows.Next() {
		var rec ScriptSmallRec
		err := rows.Scan(
			&rec.ScriptId,
			&rec.BookId,
			&rec.ChapterNum,
			//&rec.AudioFile,
			//&rec.ScriptNum,
			//&rec.UsfmStyle,
			//&rec.Person,
			//&rec.Actor,
			&rec.VerseStr)
		//&rec.ScriptText)
		if err != nil {
			log.Fatal(ctx, err, test1)
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	return results
}

func smallWordsQuery(db *sql.DB, scriptId int) []WordSmallRec {
	ctx := context.Background()
	var test1 = `SELECT w.word_id, w.word
FROM words w WHERE w.script_id = ? ORDER BY w.word_id`
	rows, err := db.Query(test1, scriptId)
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	defer rows.Close()
	var results = make([]WordSmallRec, 0, 100)
	for rows.Next() {
		var rec WordSmallRec
		err := rows.Scan(
			&rec.WordId,
			//&rec.WordSeq,
			//&rec.VerseNum,
			//&rec.Ttype,
			&rec.Word)
		if err != nil {
			log.Fatal(ctx, err, test1)
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(ctx, err, test1)
	}
	return results
}
