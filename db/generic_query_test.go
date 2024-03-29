package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGenericQuery(t *testing.T) {
	var database = `ATIWBT_USXEDIT.db`
	conn := NewDBAdapter(database)
	query := SelectType{conn.DB}
	sql1 := `SELECT book_id, chapter_num, script_begin_ts FROM scripts WHERE chapter_num=?`
	results, err := query.Select(sql1, 11)
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result[0].(string), result[1].(int), result[2].(float64))
	}
}

var test1 = `SELECT i.bible_id, i.audio_fileset_id, i.text_fileset_id, i.text_source, i.language_iso,
i.version_code, i.languge_id, s.script_id, s.book_id, s.chapter_num, s.audio_file,
s.script_num, s.usfm_style, s.person, s.actor, s.verse_str, s.script_text,
w.word_id, w.word_seq, w.verse_num, w.ttype, w.word
FROM ident i 
JOIN scripts s 
JOIN words w ON s.script_id = w.script_id
ORDER BY w.word_id`

var test1a = `SELECT s.script_id, s.book_id, s.chapter_num, s.audio_file,
s.script_num, s.usfm_style, s.person, s.actor, s.verse_str, s.script_text,
w.word_id, w.word_seq, w.verse_num, w.ttype, w.word
FROM scripts s 
JOIN words w ON s.script_id = w.script_id
ORDER BY w.word_id`

type Test1Rec struct {
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
	var db = NewDBAdapter(`BGGWFW_USXEDIT.db`)
	rows, err := db.DB.Query(test1)
	if err != nil {
		log.Fatalln(err, test1)
	}
	defer rows.Close()
	var result = make([]Test1Rec, 0, 500000)
	for rows.Next() {
		var rec Test1Rec
		err := rows.Scan(
			&rec.BibleId,
			&rec.AudioFilesetId,
			&rec.TextFilesetId,
			&rec.TextSource,
			&rec.LanguageIso,
			&rec.VersionCode,
			&rec.LanguageId,
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
			log.Fatalln(err, test1)
		}
		result = append(result, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err, test1)
	}
	fmt.Println("Count", len(result))
	fmt.Println("Elapsed", time.Since(start))
}

func TestGenericInterface(t *testing.T) {
	var start = time.Now()
	var db = NewDBAdapter(`BGGWFW_USXEDIT.db`)
	var query = SelectType{db.DB}
	records, err := query.Select(test1)
	if err != nil {
		log.Fatalln(err, test1)
	}
	fmt.Println("Query Elapsed", time.Since(start))
	var results = make([]Test1Rec, 0, 500000)
	for _, s := range records {
		var rec Test1Rec
		rec.BibleId = s[0].(string)
		rec.AudioFilesetId = s[1].(string)
		rec.TextFilesetId = s[2].(string)
		rec.TextSource = s[3].(string)
		rec.LanguageIso = s[4].(string)
		rec.VersionCode = s[5].(string)
		rec.LanguageId = s[6].(int)
		rec.ScriptId = s[7].(int)
		rec.BookId = s[8].(string)
		rec.ChapterNum = s[9].(int)
		rec.AudioFile = s[10].(string)
		rec.ScriptNum = s[11].(string)
		rec.UsfmStyle = s[12].(string)
		rec.Person = s[13].(string)
		rec.Actor = s[14].(string)
		rec.VerseStr = s[15].(string)
		rec.ScriptText = s[16].(string)
		rec.WordId = s[17].(int)
		rec.WordSeq = s[18].(int)
		rec.VerseNum = s[19].(int)
		rec.Ttype = s[20].(string)
		rec.Word = s[21].(string)
		results = append(results, rec)
	}
	fmt.Println("Count", len(results))
	fmt.Println("Total Elapsed", time.Since(start))
}
