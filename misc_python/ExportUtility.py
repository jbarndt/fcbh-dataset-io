
import os
import sys
import sqlite3
import fuzzy
import unicodedata
from DBAdapter import *
from SqliteUtility import *


class ExportUtility:


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
				ORDER BY s.script_id, w.word_id
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


	## obsolete
	def usxExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
			FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
			WHERE s.usfm_style NOT IN ('f', 'id', 'ide', 'ip', 'is', 'mt1', 'mt2', 'mt3', 'mt4', 'toc1', 'toc2', 'toc3', 'toc4', 
				'x')
			ORDER BY w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		db.close()
		return resultSet


	## obsolete
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
		isHyphen = {'\u002D', 	# hyphen-minus
				'\u2010',   	# hypen
				'\u2011',		# non-breaking hyphen
				'\u2012',		# figure dash
				'\u2013',		# en dash
				'\u2014',		# em dash
				'\u2015',		# horizontal bar
				'\uFE58',		# small em dash
				'\uFE62',		# small en dash
				'\uFE63',		# small hyphen minus
				'\uFF0D'		# fullwidth hypen-minus
				}
		name = os.path.join(os.environ.get('FCBH_DATASET_DB'), database)
		name = name.replace(".db", ".txt")
		with open(name, "w") as file:
			for row in resultSet:
				word = row[3]
				word = word.lower()
				word = unicodedata.normalize('NFD', word)
				for hyphen in isHyphen:
					word = word.replace(hyphen, '')
				result = []
				for ch in word: # remove combining diacritical marks
					if ord(ch) < 768: # x0300
						result.append(ch)
					elif ord(ch) > 879: #x036F
						result.append(ch)
				word = "".join(result)
				word = word.replace("'", '') # remove apostrophe
				word = word.replace('\uA78C', '') # remove apostrophe
				result = []
				for ch in word:
					result.append(hex(ord(ch)))
					details = " ".join(result)
				line = row[0] + ' ' + str(row[1]) +  ' ' + word + ' '  + details + '\n'
				file.write(line)


	## obsolete
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



if __name__ == "__main__":
	if len(sys.argv) < 2:
		print("Usage: python3 ExportUtility.py database")
		sys.exit(1)
	database = sys.argv[1]
	exp = ExportUtility()
	print("export ", database)
	dbpResult = exp.genericNTExport(database)
	exp.noVerseWriter(database, dbpResult)

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

