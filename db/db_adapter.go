package db

import (
	"context"
	"database/sql"
	"dataset"
	log "dataset/logger"
	"encoding/json"
	//_ "modernc.org/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

func GetDBPath(database string) string {
	var directory = os.Getenv(`FCBH_DATASET_DB`)
	if directory == `` {
		return database
	} else {
		return filepath.Join(directory, database)
	}
}

func Exists(database string) bool {
	var databasePath = GetDBPath(database)
	_, err := os.Stat(databasePath)
	return !os.IsNotExist(err)
}

func DestroyDatabase(database string) {
	var databasePath = GetDBPath(database)
	_, err := os.Stat(databasePath)
	if !os.IsNotExist(err) {
		os.Remove(databasePath)
	}
}

type DBAdapter struct {
	Ctx          context.Context
	Database     string
	DatabasePath string
	DB           *sql.DB
}

func NewDBAdapter(ctx context.Context, database string) DBAdapter {
	var databasePath = GetDBPath(database)
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(ctx, err)
	}
	execDDL(db, `PRAGMA temp_store = MEMORY;`)
	var query = `CREATE TABLE IF NOT EXISTS ident (
		dataset_id INTEGER PRIMARY KEY AUTOINCREMENT,
		bible_id TEXT NOT NULL,
		audio_fileset_id TEXT NOT NULL,
		text_fileset_id TEXT NOT NULL,
		text_source TEXT NOT NULL,
		language_iso TEXT NOT NULL,
		version_code TEXT NOT NULL,
		languge_id INTEGER,
		rolv_id INTEGER,
		alphabet TEXT,
		language_name TEXT,
		version_name TEXT) STRICT`
	execDDL(db, query)
	query = `CREATE INDEX IF NOT EXISTS ident_bible_idx ON ident (bible_id)`
	execDDL(db, query)
	query = `CREATE TABLE IF NOT EXISTS scripts (
		script_id INTEGER PRIMARY KEY AUTOINCREMENT,
		dataset_id INTEGER NOT NULL,
		book_id TEXT NOT NULL,
		chapter_num INTEGER NOT NULL,
		audio_file TEXT NOT NULL, -- questionable now that audio filesetId is in ident
		script_num TEXT NOT NULL,
		usfm_style TEXT,
		person TEXT,
		actor TEXT,  /* this should be integer. */
		verse_num INTEGER NOT NULL,
		verse_str TEXT NOT NULL, /* e.g. 6-10 7,8 6a */
		script_text TEXT NOT NULL,
		script_begin_ts REAL,
		script_end_ts REAL,
		-- script_mfcc BLOB,
		-- mfcc_rows INTEGER,
		-- mfcc_cols INTEGER,
		FOREIGN KEY(dataset_id) REFERENCES ident(dataset_id)) STRICT`
	execDDL(db, query)
	query = `CREATE UNIQUE INDEX IF NOT EXISTS scripts_idx
		ON scripts (book_id, chapter_num, script_num)`
	execDDL(db, query)
	query = `CREATE INDEX IF NOT EXISTS scripts_file_idx ON scripts (audio_file)`
	execDDL(db, query)
	query = `CREATE TABLE IF NOT EXISTS words (
		word_id INTEGER PRIMARY KEY AUTOINCREMENT,
		script_id INTEGER NOT NULL,
		word_seq INTEGER NOT NULL,
		verse_num INTEGER,
		ttype TEXT NOT NULL,
		word TEXT NOT NULL,
		word_begin_ts REAL,
		word_end_ts REAL,
		-- word_mfcc BLOB,
		-- mfcc_rows INTEGER,
		-- mfcc_cols INTEGER,
		-- mfcc_norm BLOB,
		-- mfcc_norm_rows INTEGER,
		-- mfcc_norm_cols INTEGER,
		word_enc BLOB,
		src_word_enc BLOB,
		word_multi_enc BLOB,
		src_word_multi_enc BLOB,
		FOREIGN KEY(script_id) REFERENCES scripts(script_id)) STRICT`
	execDDL(db, query)
	query = `CREATE UNIQUE INDEX IF NOT EXISTS words_idx
		ON words (script_id, word_seq)`
	execDDL(db, query)
	query = `CREATE TABLE IF NOT EXISTS script_mfcc (
		script_id INTEGER PRIMARY KEY,
		rows INTEGER NOT NULL,
		cols INTEGER NOT NULL,
		mfcc_json TEXT NOT NULL,
		FOREIGN KEY (script_id) REFERENCES scripts(script_id)) STRICT`
	execDDL(db, query)
	query = `CREATE TABLE IF NOT EXISTS word_mfcc (
		word_id INTEGER PRIMARY KEY,
		rows INTEGER NOT NULL,
		cols INTEGER NOT NULL,
		mfcc_json TEXT NOT NULL,
		FOREIGN KEY (word_id) REFERENCES words(word_id)) STRICT`
	execDDL(db, query)
	var result DBAdapter
	result.Ctx = ctx
	result.Database = database
	result.DatabasePath = databasePath
	result.DB = db
	return result
}

func (d *DBAdapter) Close() {
	err := d.DB.Close()
	if err != nil {
		log.Info(d.Ctx, err)
	}
}

//
// ident table
//

func (d *DBAdapter) InsertIdent(id Ident) dataset.Status {
	var status dataset.Status
	query := `REPLACE INTO ident(bible_id, audio_fileset_id, text_fileset_id,
		text_source, language_iso, version_code, 
		languge_id, rolv_id, alphabet, language_name, 
		version_name) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := d.DB.Prepare(query)
	defer stmt.Close()
	if err != nil {
		return log.Error(d.Ctx, 500, err, `Error while preparing Ident stmt.`)
	}
	_, err = stmt.Exec(id.BibleId, id.AudioFilesetId, id.TextFilesetId, id.TextSource, id.LanguageISO,
		id.VersionCode, id.LanguageId, id.RolvId, id.Alphabet, id.LanguageName, id.VersionName)
	if err != nil {
		return log.Error(d.Ctx, 500, err, `Error while inserting Ident.`)
	}
	return status
}

func (d *DBAdapter) InsertScripts(records []Script) dataset.Status {
	var status dataset.Status
	query := `INSERT INTO scripts(dataset_id, book_id, chapter_num, audio_file, script_num, usfm_style, 
			person, actor, verse_num, verse_str, script_text, script_begin_ts, script_end_ts) 
			VALUES (1,?,?,?,?,?,?,?,?,?,?,?,?)`
	tx, stmt := d.prepareDML(query)
	defer stmt.Close()
	for _, rec := range records {
		rec.ScriptNum = zeroFill(rec.ScriptNum, 5)
		text := strings.Join(rec.ScriptTexts, ``)
		_, err := stmt.Exec(rec.BookId, rec.ChapterNum, rec.AudioFile, rec.ScriptNum,
			rec.UsfmStyle, rec.Person, rec.Actor, rec.VerseNum, rec.VerseStr, text,
			rec.ScriptBeginTS, rec.ScriptEndTS)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error while inserting Scripts.`)
		}
	}
	status = d.commitDML(tx, query)
	return status
}

func (d *DBAdapter) SelectScalarInt(sql string) (int, dataset.Status) {
	var count int
	var status dataset.Status
	rows, err := d.DB.Query(sql)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, sql)
		return count, status
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, sql)
			return count, status
		}
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, sql)
	}
	return count, status
}

// SelectScripts is used by WordParser
func (d *DBAdapter) SelectScripts() ([]Script, dataset.Status) {
	var results []Script
	var status dataset.Status
	query := `SELECT script_id, book_id, chapter_num, verse_num, verse_str, script_text 
		FROM scripts ORDER BY script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during select scripts")
		return results, status
	}
	defer rows.Close()
	for rows.Next() {
		var rec Script
		err := rows.Scan(&rec.ScriptId, &rec.BookId, &rec.ChapterNum, &rec.VerseNum,
			&rec.VerseStr, &rec.ScriptText)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error in SelectScripts.")
			return results, status
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results, status
}

func (d *DBAdapter) SelectScriptHeadings() ([]Script, dataset.Status) {
	var result []Script
	var status dataset.Status
	query := `SELECT script_id, book_id, chapter_num, usfm_style, verse_num, verse_str, script_text 
		FROM scripts
		WHERE usfm_style IN ('para.h', 'para.mt', 'para.mt1', 'para.mt2', 'para.mt3')
		ORDER BY script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during Select Script Headings.")
		return result, status
	}
	defer rows.Close()
	for rows.Next() {
		var rec Script
		err := rows.Scan(&rec.ScriptId, &rec.BookId, &rec.ChapterNum, &rec.UsfmStyle, &rec.VerseNum,
			&rec.VerseStr, &rec.ScriptText)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error during Select Script Headings.")
			return result, status
		}
		result = append(result, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return result, status
}

// ReadNumChapters is used by match.Compare
func (d *DBAdapter) ReadNumChapters() (map[string]int, dataset.Status) {
	var results = make(map[string]int)
	var status dataset.Status
	query := `SELECT book_id, max(chapter_num) FROM scripts GROUP BY book_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, `Could not query ReadNumChapters query`)
		return results, status
	}
	type rec struct {
		bookId      string
		numChapters int
	}
	for rows.Next() {
		var tmp rec
		err := rows.Scan(&tmp.bookId, &tmp.numChapters)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, `Error reading rows in ReadNumChapters`)
			return results, status
		}
		results[tmp.bookId] = tmp.numChapters
	}
	err = rows.Err()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, `Error at end of reading rows in ReadNumChapters`)
	}
	return results, status
}

// ReadScriptsByChapter is used by Compare
func (d *DBAdapter) ReadScriptsByChapter(bookId string, chapterNum int) ([]Script, dataset.Status) {
	var results []Script
	var status dataset.Status
	sqlStmt := `SELECT book_id, chapter_num, verse_str, script_text FROM scripts 
			WHERE book_id=? AND chapter_num=?
			ORDER BY script_id`
	stmt, err := d.DB.Prepare(sqlStmt)
	defer stmt.Close()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, `Error preparing ReadScriptByChapter`)
		return results, status
	}
	rows, err := stmt.Query(bookId, chapterNum)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, `Error reading rows in ReadScriptByChapter`)
		return results, status
	}
	for rows.Next() {
		var vs Script
		err := rows.Scan(&vs.BookId, &vs.ChapterNum, &vs.VerseStr, &vs.ScriptText)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, `Error scanning in ReadScriptByChapter`)
			return results, status
		}
		results = append(results, vs)
	}
	err = rows.Err()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, `Error at end of rows in ReadingScriptByChapter`)
	}
	return results, status
}

/*
# In FileAdapter
def findChapterStart(self, book_id, chapter_num):
sql = `SELECT script_id FROM audio_scripts
WHERE book_id = ? AND chapter_num = ?
ORDER BY script_id LIMIT 1`
return self.sqlite.selectScalar(sql, [book_id, chapter_num])
*/

// SelectScriptsByBookChapter ...
func (d *DBAdapter) SelectScriptsByBookChapter(bookId string, chapter int) ([]Script, dataset.Status) {
	var results []Script
	var status dataset.Status
	var query = `SELECT script_id, script_text FROM scripts 
		WHERE book_id = ? AND chapter_num = ? ORDER BY script_id`
	rows, err := d.DB.Query(query, bookId, chapter)
	defer rows.Close()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during Select Script By Book Chapter.")
		return results, status
	}
	for rows.Next() {
		var rec = Script{BookId: bookId, ChapterNum: chapter}
		err := rows.Scan(&rec.ScriptId, &rec.ScriptText)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error during Select Script By Book Chapter.")
			return results, status
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results, status
}

func (d *DBAdapter) UpdateScriptTimestamps(scripts []Timestamp) dataset.Status {
	var status dataset.Status
	query := `UPDATE scripts SET audio_file = ?, script_begin_ts = ?,
		script_end_ts = ? WHERE script_id = ?`
	tx, stmt := d.prepareDML(query)
	defer stmt.Close()
	for _, rec := range scripts {
		_, err := stmt.Exec(rec.AudioFile, rec.BeginTS, rec.EndTS, rec.Id)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error while updating script timestamps.`)
		}
	}
	status = d.commitDML(tx, query)
	return status
}

func (d *DBAdapter) SelectScriptTimestamps(bookId string, chapter int) ([]Timestamp, dataset.Status) {
	query := `SELECT script_id, script_begin_ts, script_end_ts
		FROM scripts WHERE book_id = ? AND chapter_num = ? ORDER BY script_id`
	return d.selectTimestamps(query, bookId, chapter)
}

func (d *DBAdapter) InsertScriptMFCCS(mfccs []MFCC) dataset.Status {
	query := `REPLACE INTO script_mfcc (script_id, rows, cols, mfcc_json) VALUES (?,?,?,?)`
	return d.insertMFCCS(query, mfccs)
}

/*
def selectScriptMFCCs(self):
sql = `SELECT script_id, script_mfcc, mfcc_rows, mfcc_cols FROM audio_scripts`
resultSet = self.sqlite.select(sql)
finalSet = []
for (script_id, script_mfcc, mfcc_rows, mfcc_cols) in resultSet:
if script_mfcc != None:
mfcc_decoded = np.frombuffer(script_mfcc, dtype=np.float32)
mfcc_shaped = mfcc_decoded.reshape((mfcc_rows, mfcc_cols))
finalSet.append((script_id, mfcc_shaped))
return finalSet
*/
//
// words table
//

func (d *DBAdapter) DeleteWords() {
	execDDL(d.DB, `DELETE FROM words`)
}

func (d *DBAdapter) InsertWords(records []Word) dataset.Status {
	var status dataset.Status
	sql1 := `INSERT INTO words(script_id, word_seq, verse_num, ttype, word) VALUES (?,?,?,?,?)`
	tx, stmt := d.prepareDML(sql1)
	defer stmt.Close()
	for _, rec := range records {
		_, err := stmt.Exec(rec.ScriptId, rec.WordSeq, rec.VerseNum, rec.TType, rec.Word)
		if err != nil {
			return log.Error(d.Ctx, 500, err, "Error while inserting Words.")
		}
	}
	status = d.commitDML(tx, sql1)
	return status
}

// SelectWords is used by encode.FastText
func (d *DBAdapter) SelectWords() ([]Word, dataset.Status) {
	var results []Word
	var status dataset.Status
	var query = `SELECT word_id, ttype, word FROM words ORDER BY word_id`
	rows, err := d.DB.Query(query)
	defer rows.Close()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during Select Words.")
		return results, status
	}
	for rows.Next() {
		var rec Word
		err := rows.Scan(&rec.WordId, &rec.TType, &rec.Word)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error during Select Words.")
			return results, status
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results, status
}

// SelectWordsByBookChapter is used by Aeneas
func (d *DBAdapter) SelectWordsByBookChapter(bookId string, chapter int) ([]Word, dataset.Status) {
	var results []Word
	var status dataset.Status
	var query = `SELECT w.word_id, w.word
		FROM words w JOIN scripts s ON w.script_id = s.script_id
		WHERE w.ttype = 'W' AND s.book_id = ? AND s.chapter_num = ? ORDER BY w.word_id`
	rows, err := d.DB.Query(query, bookId, chapter)
	defer rows.Close()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during Select Words By Book Chapter.")
		return results, status
	}
	for rows.Next() {
		var rec Word
		err := rows.Scan(&rec.WordId, &rec.Word)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error during Select Words By Book Chapter.")
			return results, status
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results, status
}

/*
# In FastTextExample
def selectWords(self):
sql = `SELECT word_id, word, punct, src_word FROM audio_words`
resultSet = self.sqlite.select(sql)
return resultSet
*/
func (d *DBAdapter) UpdateWordTimestamps(words []Timestamp) dataset.Status {
	var status dataset.Status
	query := `UPDATE words SET word_begin_ts = ?, word_end_ts = ? WHERE word_id = ?`
	tx, stmt := d.prepareDML(query)
	defer stmt.Close()
	for _, rec := range words {
		_, err := stmt.Exec(rec.BeginTS, rec.EndTS, rec.Id)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error while updating word timestamps.`)
		}
	}
	status = d.commitDML(tx, query)
	return status
}

func (d *DBAdapter) SelectWordTimestamps(bookId string, chapter int) ([]Timestamp, dataset.Status) {
	query := `SELECT word_id, word_begin_ts, word_end_ts
		FROM scripts WHERE book_id = ? AND chapter_num = ? ORDER BY word_id`
	return d.selectTimestamps(query, bookId, chapter)
}

func (d *DBAdapter) selectTimestamps(query string, bookId string, chapter int) ([]Timestamp, dataset.Status) {
	var results []Timestamp
	var status dataset.Status
	rows, err := d.DB.Query(query, bookId, chapter)
	defer rows.Close()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during Select Timestamps By Book Chapter.")
		return results, status
	}
	for rows.Next() {
		var rec Timestamp
		err := rows.Scan(&rec.Id, &rec.BeginTS, &rec.EndTS)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error during Select Timestamps By Book Chapter.")
			return results, status
		}
		results = append(results, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results, status
}

func (d *DBAdapter) InsertWordMFCCS(mfccs []MFCC) dataset.Status {
	query := `REPLACE INTO word_mfcc (word_id, rows, cols, mfcc_json) VALUES (?,?,?,?)`
	return d.insertMFCCS(query, mfccs)
}

func (d *DBAdapter) insertMFCCS(query string, mfccs []MFCC) dataset.Status {
	var status dataset.Status
	tx, stmt := d.prepareDML(query)
	defer stmt.Close()
	for _, rec := range mfccs {
		mfccBytes, err := json.Marshal(rec.MFCC)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error converting MFCC to JSON`)
		}
		_, err = stmt.Exec(rec.Id, rec.Rows, rec.Cols, string(mfccBytes))
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error while inserting MFCC.`)
		}
	}
	status = d.commitDML(tx, query)
	return status
}

/*

# In MFCCExample
def selectWordMFCCs(self):
sql = `SELECT word_id, word_mfcc, mfcc_rows, mfcc_cols FROM audio_words`
resultSet = self.sqlite.select(sql, [])
finalSet = []
for (word_id, word_mfcc, mfcc_rows, mfcc_cols) in resultSet:
if word_mfcc != None:
mfcc_decoded = np.frombuffer(word_mfcc, dtype=np.float32)
mfcc_shaped = mfcc_decoded.reshape((mfcc_rows, mfcc_cols))
finalSet.append((word_id, mfcc_shaped))
return finalSet

*/

func (d *DBAdapter) UpdateWordEncodings(words []Word) dataset.Status {
	var status dataset.Status
	query := `UPDATE words SET word_enc = ? WHERE word_id = ?`
	tx, stmt := d.prepareDML(query)
	defer stmt.Close()
	for _, rec := range words {
		encBytes, err := json.Marshal(rec.WordEncoded)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error converting word enc to JSON`)
		}
		_, err = stmt.Exec(encBytes, rec.WordId)
		if err != nil {
			return log.Error(d.Ctx, 500, err, `Error while inserting word enc.`)
		}
	}
	status = d.commitDML(tx, query)
	return status
}

func execDDL(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Panic(context.Background(), err, sql)
	}
}

func zeroFill(a string, size int) string {
	num := countDigits(a)
	b := "0000000"[:size-num] + a
	return b
}

func countDigits(a string) int {
	for i, ch := range a {
		if ch < '0' || ch > '9' {
			return i
		}
	}
	return len(a)
}

func (d *DBAdapter) prepareDML(query string) (*sql.Tx, *sql.Stmt) {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Fatal(d.Ctx, err, query)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(d.Ctx, err, query)
	}
	return tx, stmt
}

func (d *DBAdapter) commitDML(tx *sql.Tx, query string) dataset.Status {
	var status dataset.Status
	err := tx.Commit()
	if err != nil {
		status = log.Error(d.Ctx, 500, err, query)
	}
	return status
}
