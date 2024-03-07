
import os
import sys
import sqlite3
import fuzzy
from DBAdapter import *
from SqliteUtility import *


class ExportAdapter:


	def genericExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
				FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
				ORDER BY w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		#self.genericWriter(database, resultSet)
		db.close()
		return resultSet


	def genericNTExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
				FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
				WHERE w.ttype = 'W'
				AND s.book_id IN ('MAT','MRK','LUK','JHN','ACT','ROM','1CO','2CO','GAL','EPH','PHP','COL',
					'1TH','2TH','1TI','2TI','TIT','PHM','HEB','JAS','1PE','2PE','1JN','2JN','3JN','JUD','REV')
				ORDER BY s.book_id, s.chapter_num, w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		#self.genericWriter(database, resultSet)
		db.close()
		return resultSet


	def genericBooksExport(self, database, books):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
				FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
				WHERE w.ttype = 'W'
				AND s.book_id IN ('""" + "','".join(books) + "') ORDER BY w.word_id"
		resultSet = db.sqlite.select(sql)
		print("results", len(resultSet))
		#self.genericWriter(database, resultSet)
		db.close()
		return resultSet


	def usxExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
			FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
			WHERE s.usfm_style NOT IN ('f', 'id', 'ide', 'ip', 'is', 'mt1', 'mt2', 'mt3', 'mt4', 'toc1', 'toc2', 'toc3', 'toc4', 
				'x')
			ORDER BY w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		#self.genericWriter(database, resultSet)
		db.close()
		return resultSet


	def genericWriter(self, database, resultSet):
		name = os.path.join(os.environ.get('FCBH_DATASET_DB'), database)
		name = name.replace(".db", ".txt")
		print("write to file", name)
		with open(name, "w") as file:
			for row in resultSet:
				word = row[3]
				word = word.replace('\u201C', '') # left quote
				word = word.replace('\u201D', '') # right quote
				word = word.replace('\u00AB', '') # <<
				word = word.replace('\u00BB', '') # >>
				word = word.replace('\u2039', '') # <
				word = word.replace('\u203A', '') # >
				word = word.replace('\u2018', '') # single left quote
				word = word.replace('\u2019', '') # single right quote
				line = row[0] + ' ' + str(row[1]) + ':' + str(row[2]) + ' ' + word + '\n'
				file.write(line)


	def noVerseWriter(self, database, resultSet):
		name = os.path.join(os.environ.get('FCBH_DATASET_DB'), database)
		name = name.replace(".db", ".txt")
		with open(name, "w") as file:
			for row in resultSet:
				word = row[3]
				word = word.lower()
				word = word.replace('\u00E9', 'e') # remove accent on e
				word = word.replace('\u00E1', 'a') # remove accent on a
				word = word.replace('\u2019', '') # remove single right quote
				word = word.replace('\u2018', '') # single left quote
				word = word.replace("'", '') # remove apostrophe
				hyphens = ['-', '\u00AD', '\u2011', '\u2012', '\u2013', '\u2014', '\u2015', '\u2043', '\u2212']
				for hyphen in hyphens:
					word = word.replace(hyphen, '') # remove hyphen
					#parts = word.split(hyphen)
					#word = parts[0]
					#if len(parts) > 1:
					#	line = row[0] + ' ' + str(row[1]) +  ' ' + parts[1] + ' '  + '\n'
					#	file.write(line)
				line = row[0] + ' ' + str(row[1]) +  ' ' + word + ' '  + '\n'
				file.write(line)


	def noVerseFuzzyWriter(self, database, resultSet):
		meta = fuzzy.DMetaphone()
		name = os.path.join(os.environ.get('FCBH_DATASET_DB'), database)
		name = name.replace(".db", ".txt")
		with open(name, "w") as file:
			for row in resultSet:
				word = row[3]
				word = word.lower()
				word = word.replace('\u2019', '') # remove single right quote
				word = word.replace("'", '') # remove apostrophe
				word = word.replace('\u2014', '') # remove hyphen
				word = word.replace('\u00E9', 'e') # remove accent on e
				word = word.replace('\u00E1', 'a') # remove accent on a
				try:
					fuzzyWord = meta(word)
					print(word, fuzzyWord)
					line = row[0] + ' ' + str(row[1]) +  ' ' + str(fuzzyWord[0]) + ' '  + '\n'
				except:
					line = row[0] + ' ' + str(row[1]) +  ' ' + str(word) + ' '  + '\n'
					print("Did not fuzzy", word)
				file.write(line)


#>>> print("DEC HEX ASC")
#... for b in bytearray(b'ABCD'):
#...     print(b, hex(b), chr(b))
#DEC HEX ASC
#65 0x41 A
#66 0x42 B
#67 0x43 C
#68 0x44 D


if __name__ == "__main__":
	exp = ExportAdapter()
	db1 = "ENGWWH_USXEDIT.db"
	db2 = "ENGWWH_WHISPER.db"
	#exp.usxExport('ZAKWYI_USX.db')
	print("export ", db1)
	dbpResult = exp.genericNTExport(db1)
	exp.noVerseWriter(db1, dbpResult)
	print("export ", db2)
	STTResult = exp.genericNTExport(db2)
	exp.noVerseWriter(db2, STTResult)
	#exp.genericExport('ZAK_MWRIGHT')
	#exp.genericExport('ENG_3_Excel')

'''
SELECT s.book_id, s.chapter_num, w.verse_num, s.usfm_style, w.word 
FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
WHERE s.book_id='MAT' AND chapter_num < 3
ORDER BY w.word_id

SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
-- WHERE s.book_id IN ('JAS') 
ORDER BY w.word_id
'''

