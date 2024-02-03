import re
import os
import csv
from DBAdapter import *
from SqliteUtility import *


class FileAdapter:

	def __init__(self, db):
		self.db = db
		self.pattern = re.compile(r"(\w+)(\W+)?$")


	def loadPoem(self, filename):
		book_id = "SHK"
		chapter_num = 1
		script_num = 0
		word_seq = 0
		verse_num = 1
		usfm_style = "q"
		person = None
		actor = None
		src_language = "ENG"
		audio_file = "../Sandeep_sample1/audio.mp3"
		with open(filename) as file:
			for line in file:
				if len(line.strip()) > 0:
					script_num += 1
					for (word_seq, word, punct) in self.parseLine(line):		
						src_word = word
						self.db.addWord(book_id, chapter_num, audio_file, script_num, usfm_style, 
							person, actor, word_seq, verse_num, word, punct, src_language, src_word)
				else:
					verse_num +=1
		self.db.insertWords()


	def loadMyDB(self, databasePath):
		script_num = 0
		word_seq = 0
		verse_num = 1
		usfm_style = "p"
		person = None
		actor = None
		src_language = None
		src_word = None
		audio_file = "../Sandeep_sample1/audio.mp3"
		srcDb = SqliteUtility(databasePath)
		resultSet = srcDb.select("SELECT reference, html FROM verses")
		for (reference, text) in resultSet:
			(book_id, chapter_num, verse_num) = reference.split(':')
			script_num = verse_num
			#print(book_id, chapter_num, verse_num, text)
			for (word_seq, word, punct) in self.parseLine(text):
				self.db.addWord(book_id, chapter_num, audio_file, script_num, usfm_style, person, actor, 
					word_seq, verse_num, word, punct, src_language, src_word)
		self.db.insertWords()
		srcDb.close()

	def loadExcel(self, filename):
		with open(filename, "r") as file:
			reader = csv.reader(file, delimiter='\t', )
			for line in reader:
				if line[0] != '' and line[1] != '':
					book_id = line[1]
					chapter_num = line[2]
					verse_num = line[3]
					person_name = line[4]
					actor_id = line[5]
					actor_name = line[6]
					script_id = line[7]
					text = line[10]
					print(line[0], "1:", line[1], "2:", line[2], "3:", line[3], 
						"4:", line[4], "5:", line[5], "6:", line[6],
						"7:", line[7], "8:", line[8], "9:", line[9], "10:", line[10])





	# This method separates punctuation
	def parseLine(self, line):
		word_seq = 0
		parts = []
		for word in line.split():
			word_seq += 1
			punct = None
			match = self.pattern.match(word)
			if match and match.group(2):
				parts.append((word_seq, match.group(1), match.group(2)))
			else:
				parts.append((word_seq, word, None))
		return parts

'''
if __name__ == "__main__":
	database = "ENG_2_WEB.db"
	if os.path.exists(database):
		os.remove(database)
	db = DBAdapter("ENG", 2, "WEB")
	file = FileAdapter(db)
	srcPath = os.environ["HOME"] + "/ShortSands/DBL/5ready/WEB.db"
	file.loadMyDB(srcPath)

#if __name__ == "__main__":
	database = "ENG_1_Sonnet.db"
	if os.path.exists(database):
		os.remove(database)
	db = DBAdapter("ENG", 1, "Sonnet")
	file = FileAdapter(db)
	file.loadPoem("../Sandeep_sample1/mplain.txt")
	resultSet = db.sqlite.select("SELECT * FROM audio_words")
	for row in resultSet:
		print(row)
'''
if __name__ == "__main__":
	database = "ENG_3_Excel.db"
	if os.path.exists(database):
		os.remove(database)
	db = DBAdapter("ENG", 3, "Excel")
	file = FileAdapter(db)
	filename = os.environ["HOME"] + "/Desktop/Mark_Scott_1_1-31-2024/excel.tsv/Script-Table 1.tsv"
	file.loadExcel(filename)

