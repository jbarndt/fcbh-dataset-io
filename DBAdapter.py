import os
import sys
import sqlite3
import numpy as np
from SqliteUtility import *


class DBAdapter:

	def __init__(self, language_iso, language_id, language_name):
		name = language_iso + "_" + str(language_id) + "_" + language_name + ".db"
		self.sqlite = SqliteUtility(name)
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
			src_language TEXT,
			src_word TEXT,
			audio_file TEXT NOT NULL,
			audio_begin_ts REAL,
			audio_end_ts REAL,
			mfccs BLOB,
			mfccs_norm BLOB,
			word_enc BLOB,
			src_word_enc BLOB,
			word_multi_enc BLOB,
			src_word_multi_enc BLOB)"""
		self.sqlite.execute(sql, [])
		sql = """CREATE UNIQUE INDEX IF NOT EXISTS audio_scripts_idx
			ON audio_words (book_id, chapter_num, script_num, word_seq)"""
		self.sqlite.execute(sql, [])


	def close(self):
		if self.sqlite != None:
			self.sqlite.close()
			self.sqlite = None


	def insertWord(self, book_id, chapter_num, script_num, word_seq, verse_num, 
		usfm_style, person, actor, word, src_language, src_word, audio_file):
		sql = """INSERT INTO audio_words(book_id, chapter_num, script_num,
			word_seq, verse_num, usfm_style, person, actor, word, 
			src_language, src_word, audio_file) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"""
		values = [book_id, chapter_num, script_num, word_seq, verse_num, 
			usfm_style, person, actor, word, src_language, src_word, audio_file]
		id = self.sqlite.executeInsert(sql, values)
		return id
		

	def updateTimestamps(self, id, audio_begin_ts, audio_end_ts):
		sql = """UPDATE audio_words SET audio_begin_ts = ?, 
			audio_end_ts = ? WHERE id = ?"""
		values = [audio_begin_ts, audio_end_ts, id]
		self.sqlite.execute(sql, values)
		# make certain word is checked here or in calling code


	def updateMFCC(self, id, mfccs):
		sql = "UPDATE audio_words SET mfccs = ? WHERE id = ?"
		values = [mfccs.tobytes(), id] # serialize mfccs
		self.sqlite.execute(sql, values)


	def updateNormalizedMFCC(self, id, mfccs_norm):
		sql = "UPDATE audio_words SET mfccs_norm = ? WHERE id = ?"
		values = [mfccs_norm.tobytes(), id] # serialize mfccs
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


	def selectWordsForFile(self, audio_file):
		sql = "SELECT id, word, src_word FROM audio_words WHERE audio_file = ?"
		resultSet = self.sqlite.select(sql, [audio_file])
		return resultSet


	def selectMFCC(self):
		sql = "SELECT id, mfccs FROM audio_words"
		resultSet = self.sqlite.select(sql, [])
		finalSet = []
		for (id, mfcc) in resultSet:
			mfcc_decoded = np.frombuffer(mfcc, dtype=double)
			finalSet.append(id, mfcc_decoded)
		return finalSet


	def selectWords(self):
		sql = "SELECT id, word, src_word FROM audio_words"
		resultSet = self.sqlite.select(sql, [])
		return resultSet


	def selectTimestamps(self, audio_file):
		sql = """SELECT id, word, audio_begin_ts, audio_end_ts
				FROM audio_words WHERE audio_file = ?"""
		resultSet = self.sqlite.select(sql, [audio_file])
		return resultSet


	def selectTensor(self):
		sql = """SELECT mfccs_norm, word_multi_enc, src_word_multi_enc 
			FROM audio_words"""
		resultSet = self.sqlite.select(sql, [])
		finalSet = []
		for (mfccs_norm, word_multi_enc, src_word_multi_enc) in resultSet:
			mfccs_decoded = np.frombuffer(mfccs_norm, dtype=double)
			word_decoded = np.frombuffer(word_multi_enc, dtype=double)
			src_word_decoded = np.frombuffer(src_word_multi_enc, dtype=double)
			finalSet.append([mfccs_decoded, word_decoded, src_word_decoded])
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
	resultSet = db.sqlite.select("SELECT id, audio_begin_ts, audio_end_ts FROM audio_words")
	for row in resultSet:
		print(row)

