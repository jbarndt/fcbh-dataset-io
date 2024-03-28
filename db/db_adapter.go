package db

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
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
	DB *sql.DB
}

func NewDBAdapter(database string) DBAdapter {
	var databasePath = GetDBPath(database)
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	var sql = `CREATE TABLE IF NOT EXISTS ident (
		bible_id TEXT NOT NULL PRIMARY KEY,
		audio_fileset_id TEXT NOT NULL,
		text_fileset_id TEXT NOT NULL,
		text_source TEXT NOT NULL,
		language_iso TEXT NOT NULL,
		version_code TEXT NOT NULL,
		languge_id INT,
		rolv_id INT,
		alphabet TEXT,
		language_name TEXT,
		version_name TEXT) STRICT`
	execDDL(db, sql)
	sql = `CREATE TABLE IF NOT EXISTS scripts (
		script_id INTEGER PRIMARY KEY AUTOINCREMENT,
		book_id TEXT NOT NULL,
		chapter_num INTEGER NOT NULL,
		audio_file TEXT NOT NULL,
		script_num TEXT NOT NULL,
		usfm_style TEXT,
		person TEXT,
		actor TEXT,  /* this should be integer. */
		verse_num INTEGER NOT NULL,
		verse_str TEXT NOT NULL, /* e.g. 6-10 7,8 6a */
		script_text TEXT NOT NULL,
		script_begin_ts REAL,
		script_end_ts REAL,
		script_mfcc BLOB,
		mfcc_rows INTEGER,
		mfcc_cols INTEGER) STRICT`
	execDDL(db, sql)
	sql = `CREATE UNIQUE INDEX IF NOT EXISTS scripts_idx
		ON scripts (book_id, chapter_num, script_num)`
	execDDL(db, sql)
	sql = `CREATE INDEX IF NOT EXISTS scripts_file_idx ON scripts (audio_file)`
	execDDL(db, sql)
	sql = `CREATE TABLE IF NOT EXISTS words (
		word_id INTEGER PRIMARY KEY AUTOINCREMENT,
		script_id INTEGER NOT NULL,
		word_seq INTEGER NOT NULL,
		verse_num INTEGER,
		ttype TEXT NOT NULL,
		word TEXT NOT NULL,
		word_begin_ts REAL,
		word_end_ts REAL,
		word_mfcc BLOB,
		mfcc_rows INTEGER,
		mfcc_cols INTEGER,
		mfcc_norm BLOB,
		mfcc_norm_rows INTEGER,
		mfcc_norm_cols INTEGER,
		word_enc BLOB,
		src_word_enc BLOB,
		word_multi_enc BLOB,
		src_word_multi_enc BLOB) STRICT`
	execDDL(db, sql)
	sql = `CREATE UNIQUE INDEX IF NOT EXISTS words_idx
		ON words (script_id, word_seq)`
	execDDL(db, sql)
	var result DBAdapter
	result.DB = db
	return result
}

func (d DBAdapter) Close() {
	err := d.DB.Close()
	if err != nil {
		log.Println(err)
	}
}

//
// ident table
//

func (d *DBAdapter) InsertIdent(bible_id string, audio_fileset_id string, text_fileset_id string,
	text_source string, language_iso string, version_code string,
	languge_id int, rolv_id int, alphabet string, language_name string,
	version_name string) {
	sql1 := `REPLACE INTO ident(bible_id, audio_fileset_id, text_fileset_id,
		text_source, language_iso, version_code, 
		languge_id, rolv_id, alphabet, language_name, 
		version_name) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	stmt, err := d.DB.Prepare(sql1)
	if err != nil {
		log.Fatal(err, sql1)
	}
	defer stmt.Close()
	_, err = stmt.Exec(bible_id, audio_fileset_id, text_fileset_id, text_source, language_iso, version_code,
		languge_id, rolv_id, alphabet, language_name, version_name)
	if err != nil {
		log.Fatal(err, sql1)
	}
}

type InsertScriptRec struct {
	BookId        string
	ChapterNum    int
	AudioFile     string
	ScriptNum     string
	UsfmStyle     string
	Person        string
	Actor         string
	VerseNum      int
	VerseStr      string
	ScriptText    []string
	ScriptBeginTs float64
	ScriptEndTs   float64
}

func (d *DBAdapter) InsertScripts(records []InsertScriptRec) {
	sql := `INSERT INTO scripts(book_id, chapter_num, audio_file, script_num, usfm_style, 
			person, actor, verse_num, verse_str, script_text, script_begin_ts, script_end_ts) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	tx, stmt := d.prepareDML(sql)
	defer stmt.Close()
	for _, rec := range records {
		rec.ScriptNum = zeroFill(rec.ScriptNum, 5)
		text := strings.Join(rec.ScriptText, ``)
		_, err := stmt.Exec(rec.BookId, rec.ChapterNum, rec.AudioFile, rec.ScriptNum,
			rec.UsfmStyle, rec.Person, rec.Actor, rec.VerseNum, rec.VerseStr, text,
			rec.ScriptBeginTs, rec.ScriptEndTs)
		if err != nil {
			log.Fatal(err, sql)
		}
	}
	d.commitDML(tx, sql)
}

func (d *DBAdapter) SelectScalarInt(sql string) int {
	rows, err := d.DB.Query(sql)
	if err != nil {
		log.Fatal(err, sql)
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			log.Fatal(err, sql)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err, sql)
	}
	return count
}

type SelectScriptsRec struct {
	ScriptId   int
	BookId     string
	ChapterNum int
	VerseNum   int
	VerseStr   string
	ScriptText string
}

// WordParser
func (d *DBAdapter) SelectScripts() []SelectScriptsRec {
	query := `SELECT script_id, book_id, chapter_num, verse_num, verse_str, script_text 
		FROM scripts ORDER BY script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		log.Fatalln(err, query)
	}
	defer rows.Close()
	var result []SelectScriptsRec
	for rows.Next() {
		var rec SelectScriptsRec
		err := rows.Scan(&rec.ScriptId, &rec.BookId, &rec.ChapterNum, &rec.VerseNum,
			&rec.VerseStr, &rec.ScriptText)
		if err != nil {
			log.Fatalln(err, query)
		}
		result = append(result, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err, query)
	}
	return result
}

type SelectScriptHeadingRec struct {
	ScriptId   int
	BookId     string
	ChapterNum int
	UsfmStyle  string
	VerseNum   int
	VerseStr   string
	ScriptText string
}

func (d *DBAdapter) SelectScriptHeadings() []SelectScriptHeadingRec {
	query := `SELECT script_id, book_id, chapter_num, usfm_style, verse_num, verse_str, script_text 
		FROM scripts
		WHERE usfm_style IN ('para.h', 'para.mt', 'para.mt1', 'para.mt2', 'para.mt3')
		ORDER BY script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		log.Fatalln(err, query)
	}
	defer rows.Close()
	var result []SelectScriptHeadingRec
	for rows.Next() {
		var rec SelectScriptHeadingRec
		err := rows.Scan(&rec.ScriptId, &rec.BookId, &rec.ChapterNum, &rec.UsfmStyle, &rec.VerseNum,
			&rec.VerseStr, &rec.ScriptText)
		if err != nil {
			log.Fatalln(err, query)
		}
		result = append(result, rec)
	}
	err = rows.Err()
	if err != nil {
		log.Fatalln(err, query)
	}
	return result
}

/*
# In FileAdapter
def findChapterStart(self, book_id, chapter_num):
sql = `SELECT script_id FROM audio_scripts
WHERE book_id = ? AND chapter_num = ?
ORDER BY script_id LIMIT 1`
return self.sqlite.selectScalar(sql, [book_id, chapter_num])

# In FileAdapter
def selectScriptsByFile(self, audio_file):
sql = `SELECT script_id, script_text FROM audio_scripts WHERE audio_file=? ORDER BY script_id`
return self.sqlite.select(sql, [audio_file])

# In FileAdapter
def addScriptTimestamp(self, script_id, script_begin_ts, script_end_ts):
self.scriptTimestampRec.append((script_begin_ts, script_end_ts, script_id))

# In FileAdapter
def updateScriptTimestamps(self):
sql = `UPDATE audio_scripts SET script_begin_ts = ?,
script_end_ts = ? WHERE script_id = ?`
self.sqlite.executeBatch(sql, self.scriptTimestampRec)
self.scriptTimestampRec = []


def selectScriptTimestamps(self, book_id, chapter_num):
sql = `SELECT script_id, script_begin_ts, script_end_ts
FROM audio_scripts WHERE book_id = ? AND chapter_num = ?`
resultSet = self.sqlite.select(sql, [book_id, chapter_num])
return resultSet


def addScriptMFCC(self, script_id, mfcc):
dims = mfcc.shape
self.scriptMfccRecs.append((mfcc.tobytes(), dims[0], dims[1], script_id))


def updateScriptMFCCs(self):
sql = `UPDATE audio_scripts SET script_mfcc = ? , mfcc_rows = ?,
mfcc_cols = ? WHERE script_id = ?`
self.sqlite.executeBatch(sql, self.scriptMfccRecs)
self.scriptMfccRecs = []


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

func (d DBAdapter) DeleteWords() {
	execDDL(d.DB, `DELETE FROM words`)
}

type InsertWordRec struct {
	ScriptId int
	WordSeq  int
	VerseNum int
	TType    string
	Word     string
}

func (d *DBAdapter) InsertWords(records []InsertWordRec) {
	sql1 := `INSERT INTO words(script_id, word_seq, verse_num, ttype, word) VALUES (?,?,?,?,?)`
	tx, stmt := d.prepareDML(sql1)
	defer stmt.Close()
	for _, rec := range records {
		_, err := stmt.Exec(rec.ScriptId, rec.WordSeq, rec.VerseNum, rec.TType, rec.Word)
		if err != nil {
			log.Fatal(err, sql1)
		}
	}
	d.commitDML(tx, sql1)
}

/*
# In AeneasExample
def selectWordsByFile(self, audio_file):
sql = `SELECT w.word_id, w.ttype, w.word
FROM audio_words w JOIN audio_scripts s ON w.script_id = s.script_id
WHERE s.audio_file=? ORDER BY w.word_id`
return self.sqlite.select(sql, [audio_file])

# In FastTextExample
def selectWords(self):
sql = `SELECT word_id, word, punct, src_word FROM audio_words`
resultSet = self.sqlite.select(sql)
return resultSet

# In Aeneas Example
def addWordTimestamp(self, word_id, word_begin_ts, word_end_ts):
self.wordTimestampRec.append((word_begin_ts, word_end_ts, word_id))

# In Aeneas Example
def updateWordTimestamps(self):
sql = `UPDATE audio_words SET word_begin_ts = ?,
word_end_ts = ? WHERE word_id = ?`
self.sqlite.executeBatch(sql, self.wordTimestampRec)
self.wordTimestampRec = []

# In MFCCExample
def selectWordTimestampsByFile(self, audio_file):
sql = `SELECT w.word_id, w.word, w.word_begin_ts, w.word_end_ts
FROM audio_words w JOIN audio_scripts s ON s.script_id = w.script_id
WHERE s.audio_file = ?`
return self.sqlite.select(sql, [audio_file])

# In MFCCExample
def addWordMFCC(self, word_id, word_mfcc):
dims = word_mfcc.shape
print(dims)
self.wordMfccRecs.append((word_mfcc.tobytes(), dims[0], dims[1], word_id))

# In MFCCExample
def updateWordMFCCs(self):
sql = `UPDATE audio_words SET word_mfcc = ? , mfcc_rows = ?,
mfcc_cols = ? WHERE word_id = ?`
self.sqlite.executeBatch(sql, self.wordMfccRecs)
self.wordMfccRecs = []

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

# In MFCCExample
def addPadWordMFCC(self, word_id, mfcc):
dims = mfcc.shape
print(dims)
self.mfccPadRecs.append((mfcc.tobytes(), dims[0], dims[1], word_id))

# In MFCCExample
def updatePadWordMFCCs(self):
sql = `UPDATE audio_words SET mfcc_norm = ?, mfcc_norm_rows = ?,
mfcc_norm_cols = ? WHERE word_id = ?`
self.sqlite.executeBatch(sql, self.mfccPadRecs)
self.mfccPadRecs = []

# In FastText Example
def addWordEncoding(self, word_id, word_enc):
self.wordEncRec.append((word_enc.tobytes(), word_id))

# In FastText Example
def updateWordEncoding(self):
sql = `UPDATE audio_words SET word_enc = ? WHERE word_id = ?`
self.sqlite.executeBatch(sql, self.wordEncRec)
self.wordEncRec = []


def addSrcWordEncoding(self, id, src_word_enc):
self.wordEncRec.append((src_word_enc.tobytes(), id))


def updateSrcWordEncoding(self):
sql = `UPDATE audio_words SET src_word_enc = ? WHERE id = ?`
self.sqlite.executeBatch(sql, self.srcWordEncRec)
self.srcWordEncRec = []


def addMultiEncoding(self, id, word_multi_enc, src_word_multi_enc):
self.multiEncRec.append((word_multi_enc.tobytes(), src_word_multi_enc.tobytes(), id))


def updateMultiEncodings(self):
sql = `UPDATE audio_words SET word_multi_enc = ?,
src_word_multi_enc = ? WHERE id = ?`
self.sqlite.executeBatch(sql, self.multiEncRec)
self.multiEncRec = []


func selectTensor() {
	var sql = `SELECT mfcc_norm, mfcc_rows, mfcc_cols, word_multi_enc,
		src_word_multi_enc FROM audio_words`
	resultSet := self.sqlite.select (sql, [])
	finalSet = []
	for (mfcc_norm, word_multi_enc, src_word_multi_enc) in resultSet:
		mfcc_decoded = np.frombuffer(mfcc_norm, dtype = np.float32)
		mfcc_shaped = mfcc_decoded.shape(mfcc_rows, mfcc_cols)
		word_decoded = np.frombuffer(word_multi_enc, dtype = np.float32)
		src_word_decoded = np.frombuffer(src_word_multi_enc, dtype = np.float32)
		finalSet.append([mfcc_shaped, word_decoded, src_word_decoded])
	return finalSet
}
*/

func execDDL(db *sql.DB, sql string) {
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatal(err, sql)
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

func (d *DBAdapter) prepareDML(sql string) (*sql.Tx, *sql.Stmt) {
	tx, err := d.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err, sql)
	}
	return tx, stmt
}

func (d *DBAdapter) commitDML(tx *sql.Tx, sql string) {
	err := tx.Commit()
	if err != nil {
		log.Fatal(err, sql)
	}
}
