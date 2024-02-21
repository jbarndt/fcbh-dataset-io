
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


	def usxExport(self, database):
		db = DBAdapter(database)
		sql = """SELECT s.book_id, s.chapter_num, w.verse_num, w.word 
			FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
			WHERE s.usfm_style NOT IN ('id', 'ide', 'mt1', 'mt2', 'mt3', 'mt4', 'f', 'x')
			ORDER BY w.word_id
		"""
		resultSet = db.sqlite.select(sql)
		self.genericWriter(database, resultSet)
		db.close()


	def genericWriter(self, database, resultSet):
		name = os.path.join(os.environ['HOME'], 'FCBH2024', database) + ".txt"
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
	exp.usxExport('ENGWEB_USX')
	exp.genericExport('WEB_1_MarkWright')
	exp.genericExport('ENG_3_Excel')

'''
SELECT s.book_id, s.chapter_num, w.verse_num, s.usfm_style, w.word 
FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id
WHERE s.book_id='MAT' AND chapter_num =2
ORDER BY w.word_id
'''

