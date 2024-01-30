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


#	def insertWord(self, book_id, chapter_num, script_num, word_seq, verse_num, 
#		usfm_style, person, actor, word, punct, src_language, src_word, audio_file):
#		sql = """INSERT INTO audio_words(book_id, chapter_num, script_num,
#			word_seq, verse_num, usfm_style, person, actor, word, punct,
#			src_language, src_word, audio_file) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)"""
#		values = [book_id, chapter_num, script_num, word_seq, verse_num, 
#			usfm_style, person, actor, word, punct, src_language, src_word, audio_file]
#		id = self.sqlite.executeInsert(sql, values)
#		return id


	def insertWord(self, book_id, chapter_num, script_num, word_seq, verse_num, 
		usfm_style, person, actor, word, punct, src_language, src_word, audio_file):
		self.insertRecs.append((book_id, chapter_num, script_num, word_seq, verse_num, 
		usfm_style, person, actor, word, punct, src_language, src_word, audio_file))


	def executeInsert(self):
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
		

	def updateTimestamps(self, id, word_begin_ts, word_end_ts):
		sql = """UPDATE audio_words SET word_begin_ts = ?, 
			word_end_ts = ? WHERE id = ?"""
		values = [word_begin_ts, word_end_ts, id]
		self.sqlite.execute(sql, values)
		# make certain word is checked here or in calling code


	def selectTimestamps(self, audio_file):
		sql = """SELECT id, word, word_begin_ts, word_end_ts
				FROM audio_words WHERE audio_file = ?"""
		resultSet = self.sqlite.select(sql, [audio_file])
		return resultSet


	def updateMFCC(self, id, mfcc):
		print("save type", type(mfcc.dtype), mfcc.shape)
		sql = """UPDATE audio_words SET mfcc = ? , mfcc_rows = ?,
			mfcc_cols = ? WHERE id = ?"""
		dims = mfcc.shape
		values = [mfcc.tobytes(), dims[0], dims[1], id] # serialize mfcc
		self.sqlite.execute(sql, values)


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


	def updateNormPaddedFCC(self, id, mfcc):
		print("save type", type(mfcc.dtype), mfcc.shape)
		sql = """UPDATE audio_words SET mfcc_norm = ?, mfcc_norm_rows = ?,
			mfcc_norm_cols = ? WHERE id = ?"""
		dims = mfcc.shape
		values = [mfcc.tobytes(), dims[0], dims[1], id] # serialize mfcc
		self.sqlite.execute(sql, values)


	def updateEncoding(self, id, word_enc):
		sql = "UPDATE audio_words SET word_enc = ? WHERE id = ?"
		values = [word_enc.tobytes(), id]
		self.sqlite.execute(sql, values)


	def updateSourceEncoding(self, id, src_word_enc):
		sql = "UPDATE audio_words SET src_word_enc = ? WHERE id = ?"
		values = [src_word_enc.tobytes(), id]
		self.sqlite.execute(sql, values)


	def updateMultiEncodings(self, id, word_multi_enc, src_word_multi_enc):
		sql = """UPDATE audio_words SET word_multi_enc = ?,
			src_word_multi_enc = ? WHERE id = ?"""
		values = [word_multi_enc.tobytes(), src_word_multi_enc.tobytes(), id]
		self.sqlite.execute(sql, values)


	def selectTensor(self):
		sql = """SELECT mfcc_norm, mfcc_rows, mfcc_cols, word_multi_enc, 
			src_word_multi_enc FROM audio_words"""
		resultSet = self.sqlite.select(sql, [])
		finalSet = []
		for (mfcc_norm, word_multi_enc, src_word_multi_enc) in resultSet:
			mfcc_decoded = np.frombuffer(mfcc_norm, dtype=np.float32)
			mfcc_shaped = mfcc_decoded.shape(mfcc_rows, mfcc_cols)
			word_decoded = np.frombuffer(word_multi_enc, dtype=np.double)
			src_word_decoded = np.frombuffer(src_word_multi_enc, dtype=np.double)
			finalSet.append([mfcc_shaped, word_decoded, src_word_decoded])
		return finalSet


if __name__ == "__main__":
	os.remove("ENG_103_English.db")
	db = DBAdapter("ENG", 103, "English")
	id1 = db.insertWord("GEN", 1, 1, 1, 1, "p", 1, 1, "In", "ENG", "In", "ENG_GEN_1.mp3")
	id2 = db.insertWord("GEN", 1, 1, 2, 1, "p", 1, 1, "the", "ENG", "the", "ENG_GEN_1.mp3")
	id3 = db.insertWord("GEN", 1, 1, 3, 1, "p", 1, 1, "beginning", "ENG", "beginning", "ENG_GEN_1.mp3")
	resultSet = db.sqlite.select("SELECT * FROM audio_words")
	for row in resultSet:
		print(row)
	db.updateTimestamps(id2, 123.45, 124.56)
	db.updateTimestamps(id3, 456.78, 456.98)
	resultSet = db.sqlite.select("SELECT id, word_begin_ts, word_end_ts FROM audio_words")
	for row in resultSet:
		print(row)

