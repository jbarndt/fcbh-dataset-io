import re
import os
import csv
from DBAdapter import *
from SqliteUtility import *


class FileAdapter:

	def __init__(self, db):
		self.db = db
		self.wordPattern = re.compile(r"(\w+)(\W+)?$")


	## for loading CSV files
	def loadExcelScripts(self, filename, audio_file_prefix):
		with open(filename, "r") as file:
			reader = csv.reader(file, delimiter='\t', )
			for line in reader:
				if line[0] != '' and line[1] != '':
					book_id = line[1]
					chapter_num = line[2]
					audio_file = audio_file_prefix + "_" + book_id + "_" + str(chapter_num).zfill(3) + "_VOX.wav"
					script_num = line[8]
					usfm_style = None
					person = line[4]
					actor = line[5]
					in_verse_num = line[3]
					if in_verse_num == "<<":
						in_verse_num = None
					script_text = line[9]
					#print("B", book_id, "C", chapter_num, "S", script_num, "P", person, "A", actor, "T", script_text)
					self.db.addScript(book_id, chapter_num, audio_file, script_num, usfm_style, person, 
					actor, in_verse_num, script_text)
			self.db.insertScripts()


	def loadWords(self):
		for (script_id, usfm_style, verse_num, script_text) in self.db.selectScripts():
			#print(script_text)
			word_seq = 0
			for word in script_text.split():
				if word[0] == '{' and word[len(word) -1] == '}':
					verse_num = word[1:len(word) -1]
					if not verse_num.isdigit(): ## A bad hack that is loosing data
						versePattern = re.compile(r"(\d+)")
						match = versePattern.match(verse_num)
						verse_num = match.group(1)
						print("Text verse num", script_id, verse_num, script_text)
				else:
					word_seq += 1
					punct = None
					match = self.wordPattern.match(word)
					if match and match.group(2):
						#print(word_seq, verse_num, match.group(1), match.group(2))
						self.db.addWord(script_id, word_seq, verse_num, match.group(1), match.group(2), None, None)
					else:
						#print(word_seq, verse_num, word)
						self.db.addWord(script_id, word_seq, verse_num, word, None, None, None)
			self.db.insertWords()


	def loadTimestamps(self, book_id, chapter_num, filename):
		first_script_id = self.db.findChapterStart(book_id, chapter_num)
		print(first_script_id)

		with open(filename, "r") as file:
			for line in file:
				(begin_ts, end_ts, line_num) = line.strip().split("\t")
				script_id = first_script_id + int(line_num) - 1
				#print(begin_ts, end_ts, line_num, script_id)
				self.db.addScriptTimestamp(script_id, begin_ts, end_ts)
			self.db.updateScriptTimestamps()


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
	database = "MZIBSM_MSCOTT.db"
	DBAdapter.destroyDatabase(database)
	db = DBAdapter(database)
	file = FileAdapter(db)
	filename = "../Mark_Scott_1_1-31-2024/excel.tsv/Script-Table 1.tsv"
	file.loadExcelScripts(filename, "N2_MZI_BSM_046")
	file.loadWords()
	filename = "../Mark_Scott_1_1-31-2024/Verse Timing File - N2_MZI_BSM_046_LUK_002_VOX.txt"
	file.loadTimestamps("LUK", 2, filename)
'''
'''
if __name__ == "__main__":
	db = DBAdapter("ZAKWYI_USX.db")
	file = FileAdapter(db)
	file.loadWords()
'''

