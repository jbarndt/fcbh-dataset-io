package output

import (
	"database/sql"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"strconv"
)

func LoadScriptStruct(d db.DBAdapter) []Script {
	var results []Script
	var status dataset.Status
	query := `SELECT scripts.script_id, book_id, chapter_num, chapter_end, audio_file, script_num, 
		usfm_style, person, actor, verse_str, verse_end, script_text, 
		script_begin_ts, script_end_ts, rows, cols, mfcc_json
		FROM scripts LEFT OUTER JOIN script_mfcc ON script_mfcc.script_id = scripts.script_id
		ORDER BY scripts.script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during select scripts")
		panic(status.Message)
	}
	defer rows.Close()
	for rows.Next() {
		var sc Script
		var mfccRows sql.NullInt64
		var mfccCols sql.NullInt64
		var mfccJson sql.NullString
		err := rows.Scan(&sc.ScriptId, &sc.BookId, &sc.ChapterNum, &sc.ChapterEnd, &sc.AudioFile,
			&sc.ScriptNum, &sc.UsfmStyle, &sc.Person, &sc.Actor, &sc.VerseStr, &sc.VerseEnd,
			&sc.ScriptText, &sc.ScriptBeginTS, &sc.ScriptEndTS, &mfccRows, &mfccCols, &mfccJson)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error in Select Scripts.")
			//return results, status
			panic(status.Message)
		}
		if mfccRows.Valid && mfccCols.Valid && mfccJson.Valid { //&& len(mfccJson.String) > 0 {
			sc.MFCCRows = int(mfccRows.Int64)
			sc.MFCCCols = int(mfccCols.Int64)
			err = json.Unmarshal([]byte(mfccJson.String), &sc.MFCC)
			if err != nil {
				panic(err)
			}
		}
		sc.Reference = FormatReference(sc.BookId, sc.ChapterNum, sc.ChapterEnd, sc.VerseStr, sc.VerseEnd)
		results = append(results, sc)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results
}

func LoadWordStruct(d db.DBAdapter) []Word {
	var results []Word
	var status dataset.Status
	query := `SELECT words.word_id, words.script_id, book_id, chapter_num, chapter_end, verse_str, 
		verse_end, words.verse_num, usfm_style, person, actor, 
		word_seq, word, word_begin_ts, word_end_ts, word_enc, rows, cols, mfcc_json
		FROM words JOIN scripts ON words.script_id = scripts.script_id
		LEFT OUTER JOIN word_mfcc ON word_mfcc.word_id = words.word_id
		WHERE ttype = 'W'
		AND book_id = 'MRK' ORDER BY words.word_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during select words")
		panic(status.Message)
	}
	defer rows.Close()
	for rows.Next() {
		var wd Word
		var wordJson string
		var mfccJson string
		err := rows.Scan(&wd.WordId, &wd.ScriptId, &wd.BookId, &wd.ChapterNum, &wd.ChapterEnd,
			&wd.VerseStr, &wd.VerseEnd, &wd.VerseNum, &wd.UsfmStyle, &wd.Person, &wd.Actor,
			&wd.WordSeq, &wd.Word, &wd.WordBeginTS, &wd.WordEndTS,
			&wordJson, &wd.MFCCRows, &wd.MFCCCols, &mfccJson)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error in Select Words.")
			//return results, status
			panic(status.Message)
		}
		err = json.Unmarshal([]byte(mfccJson), &wd.MFCC)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(wordJson), &wd.WordEnc)
		if err != nil {
			panic(err)
		}
		wd.Reference = FormatReference(wd.BookId, wd.ChapterNum, wd.ChapterEnd, wd.VerseStr, wd.VerseEnd)
		results = append(results, wd)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results
}

func FormatReference(bookId string, chapterNum int, chapterEnd int, verseStr string, verseEnd string) string {
	var result = bookId + ` ` + strconv.Itoa(chapterNum) + `:` + verseStr
	if chapterEnd != 0 && chapterNum != chapterEnd {
		result += `-` + strconv.Itoa(chapterEnd) + `:` + verseEnd
	} else if verseEnd != `` && verseStr != verseEnd {
		result += `-` + verseEnd
	}
	return result
}
