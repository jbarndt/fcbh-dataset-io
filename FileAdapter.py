import re
import os
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
						self.db.addWord(book_id, chapter_num, script_num, 
							word_seq, verse_num, usfm_style, person, actor, 
							word, punct, src_language, src_word, audio_file)
				else:
					verse_num +=1
		self.db.insertWords()
		self.db.close()


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
				self.db.addWord(book_id, chapter_num, script_num, 
							word_seq, verse_num, usfm_style, person, actor, 
							word, punct, src_language, src_word, audio_file)
		self.db.insertWords()
		self.db.close()
		srcDb.close()


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


if __name__ == "__main__":
	database = "ENG_2_WEB.db"
	if os.path.exists(database):
		os.remove(database)
	db = DBAdapter("ENG", 2, "WEB")
	file = FileAdapter(db)
	srcPath = os.environ["HOME"] + "/ShortSands/DBL/5ready/WEB.db"
	file.loadMyDB(srcPath)


'''
if __name__ == "__main__":
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
