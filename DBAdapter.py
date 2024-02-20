import os
import sys
import sqlite3
import numpy as np
from SqliteUtility import *


class DBAdapter:

	def __init__(self, database, language_id=None, language_name=None):
		#name = language_iso + "_" + str(language_id) + "_" + language_name + ".db"
		name = database + ".db"
		self.sqlite = SqliteUtility(name)
		self.scriptRecs = []
		self.scriptTimestampRec = []
		self.scriptMfccRecs = []
		self.wordRecs = []
		self.wordTimestampRec = []
		self.wordMfccRecs = []
		self.mfccPadRecs = []
		self.wordEncRec = []
		self.srcWordEncRec = []
		self.multiEncRec = []
		sql = """CREATE TABLE IF NOT EXISTS audio_scripts (
			script_id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id TEXT NOT NULL,
			chapter_num INTEGER NOT NULL,
			audio_file TEXT NOT NULL,
			script_num TEXT NOT NULL,
			-- script_sub TEXT NOT NULL,
			usfm_style TEXT,
			person TEXT,  /* should this be text or integer? */
			actor TEXT,  /* this should be integer. */
			in_verse_num INTEGER,
			script_text TEXT,
			script_begin_ts REAL,
			script_end_ts REAL,
			script_mfcc BLOB,
			mfcc_rows INTEGER,
			mfcc_cols INTEGER) STRICT"""
		self.sqlite.execute(sql)
		sql = """CREATE UNIQUE INDEX IF NOT EXISTS audio_scripts_idx
			ON audio_scripts (book_id, chapter_num, script_num)"""
		self.sqlite.execute(sql)
		sql = """CREATE INDEX IF NOT EXISTS audio_file_idx ON audio_scripts (audio_file)"""
		self.sqlite.execute(sql)
		sql = """CREATE TABLE IF NOT EXISTS audio_words (
			word_id INTEGER PRIMARY KEY AUTOINCREMENT,
			script_id INTEGER NOT NULL,
			word_seq INTEGER NOT NULL,
			verse_num INTEGER,
			word TEXT NOT NULL,
			punct TEXT,
			src_language TEXT, /* will this be replaced by script_num */
			src_word TEXT, /* will this be replaced by script_num */
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
			src_word_multi_enc BLOB) STRICT"""
		self.sqlite.execute(sql)
		sql = """CREATE UNIQUE INDEX IF NOT EXISTS audio_words_idx
			ON audio_words (script_id, word_seq)"""
		self.sqlite.execute(sql)


	def close(self):
		if self.sqlite != None:
			self.sqlite.close()
			self.sqlite = None

	#
	# audio_script table
	#

	# In FileAdapter
	def addScript(self, book_id, chapter_num, audio_file, script_num, usfm_style, 
			person, actor, in_verse_num, script_text):
		self.scriptRecs.append((book_id, chapter_num, audio_file, script_num, usfm_style, 
			person, actor, in_verse_num, script_text))

	# In FileAdapter
	def insertScripts(self):
		sql = """INSERT INTO audio_scripts(book_id, chapter_num, audio_file, 
			script_num, usfm_style, person, actor, in_verse_num, script_text) 
			VALUES (?,?,?,?,?,?,?,?,?)"""
		self.sqlite.executeBatch(sql, self.scriptRecs)
		self.scriptRecs = []

	# In FileAdapter
	def selectScripts(self):
		sql = "SELECT script_id, usfm_style, in_verse_num, script_text FROM audio_scripts"
		return self.sqlite.select(sql)

	# In FileAdapter
	def findChapterStart(self, book_id, chapter_num):
		sql = """SELECT script_id FROM audio_scripts 
				WHERE book_id = ? AND chapter_num = ?
				ORDER BY script_id LIMIT 1"""
		return self.sqlite.selectScalar(sql, [book_id, chapter_num])
				
	# In FileAdapter	
	def selectScriptsByFile(self, audio_file):
		sql = "SELECT script_id, script_text FROM audio_scripts WHERE audio_file=? ORDER BY script_id"
		return self.sqlite.select(sql, [audio_file])
				
	# In FileAdapter	
	def addScriptTimestamp(self, script_id, script_begin_ts, script_end_ts):
		self.scriptTimestampRec.append((script_begin_ts, script_end_ts, script_id))

	# In FileAdapter
	def updateScriptTimestamps(self):
		sql = """UPDATE audio_scripts SET script_begin_ts = ?, 
			script_end_ts = ? WHERE script_id = ?"""		
		self.sqlite.executeBatch(sql, self.scriptTimestampRec)
		self.scriptTimestampRec = []


	def selectScriptTimestamps(self, book_id, chapter_num):
		sql = """SELECT script_id, script_begin_ts, script_end_ts
				FROM audio_scripts WHERE book_id = ? AND chapter_num = ?"""
		resultSet = self.sqlite.select(sql, [book_id, chapter_num])
		return resultSet


	def addScriptMFCC(self, script_id, mfcc):
		dims = mfcc.shape
		self.scriptMfccRecs.append((mfcc.tobytes(), dims[0], dims[1], script_id))


	def updateScriptMFCCs(self):
		sql = """UPDATE audio_scripts SET script_mfcc = ? , mfcc_rows = ?,
			mfcc_cols = ? WHERE script_id = ?"""
		self.sqlite.executeBatch(sql, self.scriptMfccRecs)
		self.scriptMfccRecs = []


	def selectScriptMFCCs(self):
		sql = "SELECT script_id, script_mfcc, mfcc_rows, mfcc_cols FROM audio_scripts"
		resultSet = self.sqlite.select(sql)
		finalSet = []
		for (script_id, script_mfcc, mfcc_rows, mfcc_cols) in resultSet:
			if script_mfcc != None:
				mfcc_decoded = np.frombuffer(script_mfcc, dtype=np.float32)
				mfcc_shaped = mfcc_decoded.reshape((mfcc_rows, mfcc_cols))
				finalSet.append((script_id, mfcc_shaped))
		return finalSet

	#
	# audio_word table
	#

	# In FileAdapter
	def addWord(self, script_id, word_seq, verse_num, word, punct, src_language, src_word):
		self.wordRecs.append((script_id, word_seq, verse_num, word, punct, src_language, src_word))

	# In FileAdapter
	def insertWords(self):
		sql = """INSERT INTO audio_words(script_id, word_seq, verse_num, word, punct, src_language, 
			src_word) VALUES (?,?,?,?,?,?,?)"""
		self.sqlite.executeBatch(sql, self.wordRecs)
		self.wordRecs = []		

	# In AeneasExample
	def selectWordsByFile(self, audio_file):
		sql = """SELECT w.word_id, w.word, w.punct
			FROM audio_words w JOIN audio_scripts s ON w.script_id = s.script_id
			WHERE s.audio_file=? ORDER BY w.word_id"""
		return self.sqlite.select(sql, [audio_file])

	# In FastTextExample
	def selectWords(self):
		sql = "SELECT word_id, word, punct, src_word FROM audio_words"
		resultSet = self.sqlite.select(sql)
		return resultSet

	# In Aeneas Example
	def addWordTimestamp(self, word_id, word_begin_ts, word_end_ts):
		self.wordTimestampRec.append((word_begin_ts, word_end_ts, word_id))

	# In Aeneas Example
	def updateWordTimestamps(self):
		sql = """UPDATE audio_words SET word_begin_ts = ?, 
			word_end_ts = ? WHERE word_id = ?"""		
		self.sqlite.executeBatch(sql, self.wordTimestampRec)
		self.wordTimestampRec = []

	# In MFCCExample
	def selectWordTimestampsByFile(self, audio_file):
		sql = """SELECT w.word_id, w.word, w.word_begin_ts, w.word_end_ts
				FROM audio_words w JOIN audio_scripts s ON s.script_id = w.script_id
				WHERE s.audio_file = ?"""
		return self.sqlite.select(sql, [audio_file])

	# In MFCCExample
	def addWordMFCC(self, word_id, word_mfcc):
		dims = word_mfcc.shape
		print(dims)
		self.wordMfccRecs.append((word_mfcc.tobytes(), dims[0], dims[1], word_id))

	# In MFCCExample
	def updateWordMFCCs(self):
		sql = """UPDATE audio_words SET word_mfcc = ? , mfcc_rows = ?,
			mfcc_cols = ? WHERE word_id = ?"""
		self.sqlite.executeBatch(sql, self.wordMfccRecs)
		self.wordMfccRecs = []

	# In MFCCExample
	def selectWordMFCCs(self):
		sql = "SELECT word_id, word_mfcc, mfcc_rows, mfcc_cols FROM audio_words"
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
		sql = """UPDATE audio_words SET mfcc_norm = ?, mfcc_norm_rows = ?,
			mfcc_norm_cols = ? WHERE word_id = ?"""
		self.sqlite.executeBatch(sql, self.mfccPadRecs)
		self.mfccPadRecs = []

	# In FastText Example
	def addWordEncoding(self, word_id, word_enc):
		self.wordEncRec.append((word_enc.tobytes(), word_id))

	# In FastText Example
	def updateWordEncoding(self):
		sql = "UPDATE audio_words SET word_enc = ? WHERE word_id = ?"
		self.sqlite.executeBatch(sql, self.wordEncRec)
		self.wordEncRec = []	


	def addSrcWordEncoding(self, id, src_word_enc):
		self.wordEncRec.append((src_word_enc.tobytes(), id))


	def updateSrcWordEncoding(self):
		sql = "UPDATE audio_words SET src_word_enc = ? WHERE id = ?"
		self.sqlite.executeBatch(sql, self.srcWordEncRec)
		self.srcWordEncRec = []


	def addMultiEncoding(self, id, word_multi_enc, src_word_multi_enc):
		self.multiEncRec.append((word_multi_enc.tobytes(), src_word_multi_enc.tobytes(), id))


	def updateMultiEncodings(self):
		sql = """UPDATE audio_words SET word_multi_enc = ?,
			src_word_multi_enc = ? WHERE id = ?"""
		self.sqlite.executeBatch(sql, self.multiEncRec)
		self.multiEncRec = []		


	def selectTensor(self):
		sql = """SELECT mfcc_norm, mfcc_rows, mfcc_cols, word_multi_enc, 
			src_word_multi_enc FROM audio_words"""
		resultSet = self.sqlite.select(sql, [])
		finalSet = []
		for (mfcc_norm, word_multi_enc, src_word_multi_enc) in resultSet:
			mfcc_decoded = np.frombuffer(mfcc_norm, dtype=np.float32)
			mfcc_shaped = mfcc_decoded.shape(mfcc_rows, mfcc_cols)
			word_decoded = np.frombuffer(word_multi_enc, dtype=np.float32)
			src_word_decoded = np.frombuffer(src_word_multi_enc, dtype=np.float32)
			finalSet.append([mfcc_shaped, word_decoded, src_word_decoded])
		return finalSet


if __name__ == "__main__":
	database = "ENG_103_English.db"
	if os.path.exists(database):
		os.remove(database)
	db = DBAdapter("ENG_103_English", 103, "English")

	print("* Expect 3 lines of script records")
	db.addScript("GEN", 1, "ENG_GEN_1.mp3", 1, "p", 1, 1, 1, "In the beginning darkness")
	db.addScript("GEN", 1, "ENG_GEN_1.mp3", 2, "p", 1, 1, 1, "Let there be light")
	db.addScript("GEN", 1, "ENG_GEN_1.mp3", 3, "p", 1, 1, 2, "And there was light")
	db.insertScripts()
	resultSet = db.selectScriptsByFile("ENG_GEN_1.mp3")
	for (script_id, script_text) in resultSet:
		print(script_id, script_text)

	print("* Expect 3 lines of timestamps")
	db.addScriptTimestamp(1, 0.0, 123.45)
	db.addScriptTimestamp(2, 123.45, 124.56)
	db.addScriptTimestamp(3, 456.78, 456.98)
	db.updateScriptTimestamps()
	resultSet = db.selectScriptTimestamps("GEN", 1)
	for (script_id, script_begin_ts, script_end_ts) in resultSet:
		print(script_id, script_begin_ts, script_end_ts)

	print("* Expect 1 line of MFCC records")
	mfcc = np.array([[1, 2, 3], [4, 5, 6]], np.int32)
	db.addScriptMFCC(1, mfcc)
	db.updateScriptMFCCs()
	resultSet = db.selectScriptMFCCs()
	for (script_id, mfcc) in resultSet:
		print(script_id, mfcc)

	print("* Expect 3 lines of word records")
	db.addWord(1, 1, 1, "In", None, "ENG", "In")
	db.addWord(1, 2, 1, "the", None, "ENG", "the")
	db.addWord(1, 3, 1, "beginning", None, "ENG", "beginning")
	db.insertWords()
	resultSet = db.selectWordsByFile("ENG_GEN_1.mp3")
	for (word_id, word, punct) in resultSet:
		print(word_id, word, punct)

	print("* Expect 3 lines of word timestamps")
	db.addWordTimestamp(1, 0, 123.45)
	db.addWordTimestamp(2, 123.45, 124.56)
	db.addWordTimestamp(3, 456.78, 456.98)
	db.updateWordTimestamps()
	resultSet = db.selectWordTimestampsByFile("ENG_GEN_1.mp3")
	for row in resultSet:
		print(row)

	print("* Expect 1 line of word MFCC")
	word_mfcc = np.array([[1, 2, 3], [4, 5, 6]], np.int32)
	db.addWordMFCC(1, word_mfcc)
	db.updateWordMFCCs()
	resultSet = db.selectWordMFCCs()
	for row in resultSet:
		print(row)

#SELECT w.script_id, w.word_id, w.word_seq, w.word, w.punct, w.src_word 
#FROM audio_words w JOIN audio_scripts s ON w.script_id = s.script_id
#WHERE s.book_id='GEN' AND s.chapter_num=1

