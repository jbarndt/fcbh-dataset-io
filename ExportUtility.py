
import os
import sys
import sqlite3
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
		self.genericWriter(database, resultSet)
		db.close()


	def genericBooksExport(self, database, books):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
				FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
				WHERE s.book_id IN ('""" + "','".join(books) + "') ORDER BY w.word_id"
		resultSet = db.sqlite.select(sql)
		self.genericWriter(database, resultSet)
		db.close()


	def usxExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
			FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
			WHERE s.usfm_style NOT IN ('f', 'id', 'ide', 'ip', 'is', 'mt1', 'mt2', 'mt3', 'mt4', 'toc1', 'toc2', 'toc3', 'toc4', 
				'x')
			ORDER BY w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		self.genericWriter(database, resultSet)
		db.close()


	def genericWriter(self, database, resultSet):
		name = os.path.join(os.environ.get('FCBH_DATASET_DB'), database)
		name = name.replace(".db", ".txt")
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


if __name__ == "__main__":
	exp = ExportAdapter()
	exp.usxExport('ZAKWYI_USX.db')
	exp.genericBooksExport('ZAK_MWRIGHT.db', ['JAS'])
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

