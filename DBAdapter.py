import os
import sys
import sqlite3
import numpy as np
from SqliteUtility import *


class DBAdapter:

	def __init__(self, language_iso, language_id, language_name):
		name = language_iso + "_" + str(language_id) + "_" + language_name + ".db"
		self.sqlite = SqliteUtility(name)
		self.insertRecs = []
		self.timestampRec = []
		self.mfccRecs = []
		self.mfccPadRecs = []
		self.wordEncRec = []
		self.srcWordEncRec = []
		self.multiEncRec = []
		sql = """CREATE TABLE IF NOT EXISTS audio_words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id TEXT NOT NULL,
			chapter_num INTEGER NOT NULL,
			script_num INTEGER NOT NULL,
			word_seq INTEGER NOT NULL,
			verse_num INTEGER,
			usfm_style TEXT,
			person INTEGER,
			actor INTEGER,
			word TEXT NOT NULL,
			punct TEXT,
			src_language TEXT, /* will this be replaced by script_num */
			src_word TEXT, /* will this be replaced by script_num */
			audio_file TEXT NOT NULL,
			word_begin_ts REAL,
			word_end_ts REAL,
			mfcc BLOB,
			mfcc_rows INTEGER,
			mfcc_cols INTEGER,
			mfcc_norm BLOB,
			mfcc_norm_rows INTEGER,
			mfcc_norm_cols INTEGER,
			word_enc BLOB,
			src_word_enc BLOB,
			word_multi_enc BLOB,
			src_word_multi_enc BLOB)"""
		self.sqlite.execute(sql)
		sql = """CREATE UNIQUE INDEX IF NOT EXISTS audio_scripts_idx
			ON audio_words (book_id, chapter_num, script_num, word_seq)"""
		self.sqlite.execute(sql)


	def close(self):
		if self.sqlite != None:
			self.sqlite.close()
			self.sqlite = None



	def addWord(self, book_id, chapter_num, script_num, word_seq, verse_num, 
		usfm_style, person, actor, word, punct, src_language, src_word, audio_file):
		self.insertRecs.append((book_id, chapter_num, script_num, word_seq, verse_num, 
		usfm_style, person, actor, word, punct, src_language, src_word, audio_file))


	def insertWords(self):
		sql = """INSERT INTO audio_words(book_id, chapter_num, script_num,
			word_seq, verse_num, usfm_style, person, actor, word, punct,
			src_language, src_word, audio_file) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"""
		self.sqlite.executeBatch(sql, self.insertRecs)
		self.insertRecs = []		


	def selectWordsForFile(self, audio_file):
		sql = "SELECT id, word, src_word FROM audio_words WHERE audio_file = ?"
		resultSet = self.sqlite.select(sql, [audio_file])
		return resultSet


	def selectWords(self):
		sql = "SELECT id, word, src_word FROM audio_words"
		resultSet = self.sqlite.select(sql, [])
		return resultSet

	def selectScript(self):
		sql = """SELECT id, book_id, chapter_num, script_num, word_seq, 
			verse_num, usfm_style, person, word, punct FROM audio_words"""
		resultSet = self.sqlite.select(sql)
		return resultSet


	def addTimestamp(self, id, word_begin_ts, word_end_ts):
		self.timestampRec.append((word_begin_ts, word_end_ts, id))


	def updateTimestamps(self):
		sql = """UPDATE audio_words SET word_begin_ts = ?, 
			word_end_ts = ? WHERE id = ?"""		
		self.sqlite.executeBatch(sql, self.timestampRec)
		self.timestampRec = []


	def selectTimestamps(self, audio_file):
		sql = """SELECT id, word, word_begin_ts, word_end_ts
				FROM audio_words WHERE audio_file = ?"""
		resultSet = self.sqlite.select(sql, [audio_file])
		return resultSet


	def addMFCC(self, id, mfcc):
		#print("save type", type(mfcc.dtype), mfcc.shape)
		dims = mfcc.shape
		self.mfccRecs.append((mfcc.tobytes(), dims[0], dims[1], id))


	def updateMFCCs(self):
		sql = """UPDATE audio_words SET mfcc = ? , mfcc_rows = ?,
			mfcc_cols = ? WHERE id = ?"""
		self.sqlite.executeBatch(sql, self.mfccRecs)
		self.mfccRecs = []


	def selectMFCC(self):
		sql = "SELECT id, mfcc, mfcc_rows, mfcc_cols FROM audio_words"
		resultSet = self.sqlite.select(sql, [])
		finalSet = []
		for (id, mfcc, mfcc_rows, mfcc_cols) in resultSet:
			mfcc_decoded = np.frombuffer(mfcc, dtype=np.float32)
			mfcc_shaped = mfcc_decoded.reshape((mfcc_rows, mfcc_cols))
			print(mfcc_decoded.shape, mfcc_shaped.shape)
			finalSet.append((id, mfcc_shaped))
		return finalSet


	def addPadMFCC(self, id, mfcc):
		dims = mfcc.shape
		self.mfccPadRecs.append((mfcc.tobytes(), dims[0], dims[1], id))	


	def updatePadMFCCs(self):
		sql = """UPDATE audio_words SET mfcc_norm = ?, mfcc_norm_rows = ?,
			mfcc_norm_cols = ? WHERE id = ?"""
		self.sqlite.executeBatch(sql, self.mfccPadRecs)
		self.mfccPadRecs = []


	def addWordEncoding(self, id, word_enc):
		self.wordEncRec.append((word_enc.tobytes(), id))


	def updateWordEncoding(self):
		sql = "UPDATE audio_words SET word_enc = ? WHERE id = ?"
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
	os.remove("ENG_103_English.db")
	db = DBAdapter("ENG", 103, "English")
	id1 = db.addWord("GEN", 1, 1, 1, 1, "p", 1, 1, "In", None, "ENG", "In", "ENG_GEN_1.mp3")
	id2 = db.addWord("GEN", 1, 1, 2, 1, "p", 1, 1, "the", None, "ENG", "the", "ENG_GEN_1.mp3")
	id3 = db.addWord("GEN", 1, 1, 3, 1, "p", 1, 1, "beginning", ",", "ENG", "beginning", "ENG_GEN_1.mp3")
	db.insertWords()
	resultSet = db.sqlite.select("SELECT * FROM audio_words")
	for row in resultSet:
		print(row)
	db.addTimestamp(id2, 123.45, 124.56)
	db.addTimestamp(id3, 456.78, 456.98)
	db.updateTimestamps()
	resultSet = db.sqlite.select("SELECT id, word_begin_ts, word_end_ts FROM audio_words")
	for row in resultSet:
		print(row)

