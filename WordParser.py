import re
import os
from DBAdapter import *

# This program reads text scripts from the audio_script table and produces the audio_words table.
# It outputs the text as token types and tokens.  The following are the token types:
# W Word: this is whole words, including hypenated words
# P Punctuation: this is single punctuation characters
# S Whitespace: this is exact whitespace, and can be multiple characters.
# V Verse number: {n} or {n}_{n} this is the exact character string as shown.
# N Number (was considered, but not done)
# While languages might differ, this program needs a consistent way to decide what is a word.
# It assumes that whitespace is always a word delimiter, but a single character of punctuation
# is not a word delimiter.  But, it assumes that multiple characters of punctuation within a word
# is word delimiter.


class WordParser:

	def __init__(self, db):
		self.db = db
		self.word_seq = 0
		self.lastScriptId = None


	def parse(self):
		BEGIN = 1
		SPACE = 2
		WORD = 3 
		WORDPUNCT = 4
		VERSENUM = 5
		INVERSENUM = 6
		ENDVERSENUM = 7
		NEXTVERSENUM = 8
		label = ["", "BEGIN", "SPACE", "WORD", "WORDPUNCT", "VERSENUM", "INVERSENUM", "ENDVERSENUM", 
			"NEXTVERSENUM"] 
		for (script_id, usfm_style, verse_num, script_text) in self.db.selectScripts():
			print(script_text)
			term = None
			punct = None
			verseeStr = None
			state = BEGIN
			for token in re.split(r"([\W])", script_text):
				if len(token) > 0:
					print(label[state], "now token: ", token)
					if state == BEGIN:
						if token.isspace():
							term = token
							state = SPACE
						elif token.isalnum(): 
							term = token
							state = WORD   
						elif token == '{':
							term = token
							state = VERSENUM
						else: # token.ispunct()
							self.addWord(script_id, verse_num, 'P', token)
							term = None # redundant
							state = BEGIN
					elif state == SPACE:
						if token.isspace():
							term += token
						elif token.isalnum():
							self.addWord(script_id, verse_num, 'S', term)
							term = token
							state = WORD
						elif token == '{':
							self.addWord(script_id, verse_num, 'S', term)
							term = token
							state = VERSENUM
						else: # token.ispunct()
							self.addWord(script_id, verse_num, 'S', term)
							self.addWord(script_id, verse_num, 'P', token)
							term = None
							state = BEGIN
					elif state == WORD:
						if token.isspace():
							self.addWord(script_id, verse_num, 'W', term)
							term = token 
							state = SPACE
						elif token.isalnum():
							self.logError("space or punct", token)
						elif token == '{':
							self.addWord(script_id, verse_num, 'W', term)
							term = token
							state = VERSENUM
						else: # token.ispunct()
							punct = token 
							state = WORDPUNCT
					elif state == WORDPUNCT:
						if token.isspace():
							self.addWord(script_id, verse_num, 'W', term)
							self.addWord(script_id, verse_num, 'P', punct)
							punct = None
							term = token 
							state = SPACE
						elif token.isalnum():
							term += punct + token 
							punct = None 
							state = WORD
						else: # token.ispunct()
							self.addWord(script_id, verse_num, 'W', term)
							self.addWord(script_id, verse_num, 'P', punct)
							self.addWord(script_id, verse_num, 'P', token)
							term = None
							state = BEGIN
					elif state == VERSENUM:
						if token.isdigit():
							term += token 
							verseStr = token
							state = INVERSENUM
						else:
							self.logError("number", token)
					elif state == INVERSENUM:
						if token.isdigit():
							term += token
							verseStr += token
						elif token == '}':
							term += token
							verse_num = int(verseStr)
							state = ENDVERSENUM
						else:
							self.logError("number or }", token)
					elif state == ENDVERSENUM:
						if token == '_':
							term += token 
							state = NEXTVERSENUM
						elif token.isspace():
							self.addWord(script_id, verse_num, 'V', term)
							term = token
							state = SPACE
						elif token.isalnum():
							self.addWord(script_id, verse_num, 'V', term)
							term = token 
							state = WORD 
						else: # token.ispunct()
							print("output verse", term)
							self.addWord(script_id, verse_num, 'V', term)
							self.addWord(script_id, verse_num, 'P', token)
							term = None 
							state = BEGIN
					elif state == NEXTVERSENUM:
						if token == '{':
							term += token
							state = INVERSENUM
						else:
							self.logError("{", token)

			if term != None and len(term) > 0:
				if term.isspace():
					self.addWord(script_id, verse_num, 'S', term)
				elif term.isalnum():
					self.addWord(script_id, verse_num, 'W', term)
		self.db.insertWords()


	def addWord(self, script_id, verse_num, ttype, text):
		if self.lastScriptId != script_id:
			self.lastScriptId = script_id
			self.word_seq = 0
		self.word_seq += 1
		if ttype == None or text == None:# or verse_num == None:
			print(script_id, verse_num, ttype, text)
			sys.exit(0)
		print("seq: ", self.word_seq, " versee: ", verse_num, " type: ", ttype, " text: ", text)
		self.db.addWord(script_id, self.word_seq, verse_num, ttype, text)


	def logError(self, expected, actual):
		print("Expected: ", expected, ", but found: ", actual)
		sys.exit(1)


	def format(self, filename):
		with open(filename, "w") as file:
			sql = """SELECT s.script_id, s.book_id, s.chapter_num, w.verse_num, w.word_seq, w.ttype, w.word 
				FROM audio_scripts s JOIN audio_words w ON s.script_id = w.script_id ORDER BY w.word_id"""
			for (script_id, book_id, chapter_num, verse_num, word_seq, ttype, word) in db.sqlite.select(sql):
				file.write(word)


if __name__ == "__main__":
	db = DBAdapter("ZAK_MWRIGHT.db")
	db.deleteWords()
	word = WordParser(db)
	word.parse()
	word.format("ZAK_MWRIGHT_WORDS.txt")
	with open("ZAK_MWRIGHT_SCRIPT.txt", "w") as file:
		sql = "SELECT book_id, chapter_num, in_verse_num, script_text FROM audio_scripts ORDER BY script_id"
		for (book_id, chapter_num, in_verse_num, script_text) in db.sqlite.select(sql):
			file.write(script_text)





