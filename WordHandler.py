import re
import os
import csv
from DBAdapter import *
from SqliteUtility import *


class WordHandler:

	def __init__(self, db):
		self.db = db


	def wordParser(self):
		BEGIN = 1
		VERSENUM = 2
		ENDVERSENUM = 3
		PREPUNCT = 4
		WORD = 5
		POSTPUNCT = 6
		label = ["", "BEGIN", "VERSENUM", "ENDVERSENUM", "PREPUNCT", "WORD", "POSTPUNCT"]
		for (script_id, usfm_style, verse_num, script_text) in self.db.selectScripts():
			print(script_text)
			word_seq = 0
			prePunct = trueWord = postPunct = None
			state = BEGIN
			for word in re.split(r"([\W])", script_text):
				print("HAS", label[state], word)
				if len(word) > 0:
					if state == BEGIN:
						if word == '{':
							state = VERSENUM
						elif word.isspace():
							state = BEGIN
						elif word.isalnum():
							state = WORD
							trueWord = word
						else:
							state = PREPUNCT
							prePunct = word
					elif state == VERSENUM:
						if word.isdigit():
							state = ENDVERSENUM
							verse_num = word 
						else:
							self.logError("number", word)
					elif state == ENDVERSENUM:
						if word == '}':
							state = BEGIN
						else:
							self.logError("}", word)
					elif state == PREPUNCT:
						if word.isalnum():
							state = WORD 
							trueWord = word
						elif word.isspace():
							state = BEGIN 
							word_seq += 1
							#self.db.addWord(script_id, word_seq, verse_num, prePunct, trueWord, postPunct, None, None)
							print("pre:", prePunct, "word:", trueWord, "post:", postPunct)
							prePunct = trueWord = postPunct = None
						else:
							prePunct += word
							#self.logError("a word", word)
					elif state == WORD:
						if word.isspace():
							state = BEGIN
							word_seq += 1
							#self.db.addWord(script_id, word_seq, verse_num, prePunct, trueWord, postPunct, None, None)
							print("pre:", prePunct, "word:", trueWord, "post:", postPunct)
							prePunct = trueWord = postPunct = None
						elif not word.isalnum():
							state = POSTPUNCT
							postPunct = word
						else:
							self.logError("punct or whitespace", word)
					elif state == POSTPUNCT:
						if word.isspace():
							state = BEGIN 
							word_seq += 1
							#self.db.addWord(script_id, word_seq, verse_num, prePunct, trueWord, postPunct, None, None)
							print("pre:", prePunct, "word:", trueWord, "post:", postPunct)
							prePunct = trueWord = postPunct = None
						elif not word.isalnum():
							postPunct += word
						else:
							#self.logError("whitespace", word)
							state = WORD
							trueWord += postPunct + word
							postPunct = None
					else:
						self.logError("Unknown state", state)

	def logError(self, expected, actual):
		print("Expected: ", expected, ", but found: ", actual)
		sys.exit(1)


if __name__ == "__main__":
	db = DBAdapter("ZAK_MWRIGHT.db")
	word = WordHandler(db)
	word.wordParser()



