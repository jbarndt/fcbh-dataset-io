
import os
import sys
import urllib.request
import json
from DBAdapter import *
from FileAdapter import *

# This program loads a plain text fileset into a sqlite database
# it loads it from the downloads directory.

class DBPTextAdapter:


	def __init__(self, db):
		self.db = db


	def insertIdent(self, bibleId):
		iso = bibleId[:3]
		version = bibleId[3:]
		self.db.insertIdent(bibleId, iso, version, "DBPTEXT", None, None, None, None, None)


	def processDirectory(self, bibleId):
		filesetId = bibleId + "N_ET"
		directory = os.path.join(os.environ['FCBH_DATASET_FILES'], bibleId)
		self.processFile(directory, bibleId + "O_ET.json")
		self.processFile(directory, bibleId + "N_ET.json")


	def processFile(self, directory, filename):
		scriptNum = 0
		filePath = os.path.join(directory, filename)
		lastBookId = None
		try:
			with open(filePath, "rb") as file:
				content = file.read()
				print("Read", filename, len(content), "bytes")
				verses = json.loads(content.decode('utf-8'))
				print("num verses", len(verses))
				for vs in verses:
					scriptNum += 1
					bookId = vs['book_id']
					if lastBookId != bookId:
						print(bookId)
						lastBookId = bookId
					chapter = vs['chapter']
					verseNum = vs['verse_start']
					text = vs['verse_text']
					self.db.addScript(bookId, chapter, filename, scriptNum, None, None, None, verseNum, text)
		except Exception as e:
			print("Error", e)
		self.db.insertScripts()


if __name__ == "__main__":
	if len(sys.argv) < 2:
		print("Usage: DBPTextAdapter  bibleId")
		sys.exit(1)
	bibleId = sys.argv[1]
	database = bibleId + "_DBPTEXT.db"
	DBAdapter.destroyDatabase(database)
	db = DBAdapter(database)
	text = DBPTextAdapter(db)
	text.insertIdent(bibleId)
	text.processDirectory(bibleId)
	file = FileAdapter(db)
	file.loadWords()
	db.close()


