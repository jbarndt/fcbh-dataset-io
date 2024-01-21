#import os
#import io
import re
from DBAdapter import *


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
		pattern = re.compile(r"(\w+)(\W+)?$")
		with open(filename) as file:
			for line in file:
				if len(line.strip()) > 0:
					script_num += 1
					word_seq = 0
					words = self.parseLine(line)
					for word in words:		
						src_word = word
						word_seq += 1
						db.insertWord(book_id, chapter_num, script_num, 
							word_seq, verse_num, usfm_style, person, actor, 
							word, src_language, src_word, audio_file)
				else:
					verse_num +=1


	def parseLine(self, line):
		parts = []
		for word in line.split():
			punct = None
			match = self.pattern.match(word)
			if match and match.group(2):
				parts.append(match.group(1))
				parts.append(match.group(2))
			else:
				parts.append(word)
		return parts


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

