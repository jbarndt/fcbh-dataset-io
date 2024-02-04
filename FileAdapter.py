import re
import os
import csv
from DBAdapter import *
from SqliteUtility import *


class FileAdapter:

	def __init__(self, db):
		self.db = db
		self.wordPattern = re.compile(r"(\w+)(\W+)?$")
		self.numPattern = re.compile(r"(\d+)(\D+)")

	## Obsolete
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

	## Obsolete
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


	def loadExcelScripts(self, filename, audio_file_prefix):
		with open(filename, "r") as file:
			reader = csv.reader(file, delimiter='\t', )
			for line in reader:
				if line[0] != '' and line[1] != '':
					book_id = line[1]
					chapter_num = line[2]
					audio_file = audio_file_prefix + "_" + book_id + "_" + chapter_num + ".vox"
					match = re.match(self.numPattern, line[8])
					if match:
						script_num = match.group(1)
						script_sub = match.group(2)
					else:
						script_num = line[8]
						script_sub = ''
					#print(line[8], script_num, script_sub, script_sub == '')
					usfm_style = None
					person = line[4]
					actor = line[5]
					in_verse_num = line[3]
					if in_verse_num == "<<":
						in_verse_num = None
					script_text = line[10]
					#print("B", book_id, "C", chapter_num, "S", script_num, "P", person, "A", actor, "T", script_text)
					self.db.addScript(book_id, chapter_num, audio_file, script_num, script_sub, usfm_style, person, 
					actor, in_verse_num, script_text)
			self.db.insertScripts()


	def loadWords(self):
		for (script_id, usfm_style, verse_num, script_text) in self.db.selectScripts():
			#print(script_text)
			word_seq = 0
			for word in script_text.split():
				if word[0] == '{' and word[len(word) -1] == '}':
					verse_num = word[1:len(word) -1]
				else:
					word_seq += 1
					punct = None
					match = self.wordPattern.match(word)
					if match and match.group(2):
						#print(word_seq, verse_num, match.group(1), match.group(2))
						db.addWord(script_id, word_seq, verse_num, match.group(1), match.group(2), None, None)
					else:
						#print(word_seq, verse_num, word)
						db.addWord(script_id, word_seq, verse_num, word, None, None, None)
			db.insertWords()


	def loadTimestamps(self, book_id, chapter_num, filename):
		first_script_id = self.db.findChapterStart(book_id, chapter_num)
		with open(filename, "r") as file:
			for line in file:
				(begin_ts, end_ts, line_num) = line.strip().split("\t")
				script_id = first_script_id + int(line_num) - 2
				#print(begin_ts, end_ts, line_num, script_id)
				db.addScriptTimestamp(script_id, begin_ts, end_ts)
			db.updateScriptTimestamps()



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
	if os.path.exists(database):, 
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
	file.loadExcelScripts(filename, "N2_MZI_BSM_046")
	file.loadWords()
	filename = os.environ["HOME"] + "/Desktop/Mark_Scott_1_1-31-2024/Verse Timing File - N2_MZI_BSM_046_LUK_002_VOX.txt"
	file.loadTimestamps("LUK", 2, filename)

