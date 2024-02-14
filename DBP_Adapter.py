
import os
import sys
import urllib.request
import json
from DBAdapter import *
from FileAdapter import *

class DBP_Adapter:


	def process(self):
		(fileset_id, book_id, chapter_num) = self.arguments()
		(fileType, inCloud) = self.findType(fileset_id)
		print(fileset_id, book_id, chapter_num, fileType, inCloud)
		if inCloud:
			directory = os.environ["HOME"] + "/Downloads"
			self.downloadFile(directory, fileType, fileset_id, book_id, chapter_num)
		else:
			content1 = self.downloadJson(fileset_id)
			print(content1["data"][0])
			print()
			print(content1["meta"])
			db = self.createDatabase(fileset_id)
			self.loadAudioScript(db, content1["data"])
			file = FileAdapter(db)
			file.loadWords()


	def arguments(self):
		print(sys.argv)
		if len(sys.argv) == 1:
			print("Usage: python3 DBP_Adapter.py  fileset_id  [book_id]  [chapter_num]")
			sys.exit(1
				)
		fileset_id = sys.argv[1] if len(sys.argv) > 1 else None
		book_id = sys.argv[2] if len(sys.argv) > 2 else None
		chapter = sys.argv[3] if len(sys.argv) > 3 else '0'
		if chapter.isdigit():
			chapter_num = int(chapter)
		else:
			print("The third parameter is chapter_num and must be a number.")
			sys.exit(1)
		return(fileset_id, book_id, chapter_num)


	def findType(self, fileset_id):
		parts = fileset_id.split("-")
		if len(parts) > 1:
			if parts[1] == "usx":
				return("usx", True)
			elif parts[1] == "opus16":
				return("opus", True)
			elif parts[1] == "json":
				return("json", True)
			elif parts[1] == "mp3":
				return("mp3", True)
			else:
				print("fileset_id has an unknown type", fileset_id)
				sys.exit(1)
		else:
			filesetType = fileset_id[-2:]
			if filesetType == "DA":
				return("mp3", True)
			elif filesetType == "ET":
				return("json", False)
			elif filesetType == "SA":
				print("Unable to process fileset of type SA")
				sys.exit(1)
			else:
				print("fileset_id has an unknown type", fileset_id)
				sys.exit(1)


	def downloadJson(self, fileset_id, page = None):
		url = "https://4.dbt.io/api/download/" + fileset_id + "?v=4&key=b4715786-9b8e-4fbe-a9b9-ff448449b81b"
		if page != None:
			url += "&page=" + str(page)
		print(url)
		content = self.jsonRequest(fileset_id, url)
		return content 


	def createDatabase(self, fileset_id):
		isoCode = fileset_id[:3]
		versionCode = fileset_id[3:6]
		database = isoCode + "_1_" + versionCode + ".db"
		if os.path.exists(database):
			os.remove(database)
		db = DBAdapter(isoCode, 1, versionCode)
		return db


	def loadAudioScript(self, db, elements):
		for index, ele in enumerate(elements):
			bookSeq = bookSeqMap[ele['book_id']]
			elements[index]['book_seq'] = bookSeq
		sorted_elements = sorted(elements, key=lambda x: (x['book_seq'], x['chapter'], x['verse_start']))
		script_num = 0
		for rec in sorted_elements:
			book_id = rec['book_id']
			chapter_num = rec['chapter']
			audio_file = "None"
			script_num += 1
			usfm_style = None 
			person = None  
			actor = None
			in_verse_num = rec['verse_start']
			script_text = rec['verse_text']
			db.addScript(book_id, chapter_num, audio_file, script_num, usfm_style, 
				person, actor, in_verse_num, script_text)
		db.insertScripts()


	def downloadFile(self, directory, fileType, fileset_id, book_id, chapter_num):
		url = "https://4.dbt.io/api/download/" + fileset_id + "?v=4&key=b4715786-9b8e-4fbe-a9b9-ff448449b81b"
		content = self.jsonRequest(fileset_id, url)
		## need to look at content['meta'] for pagination
		items = content['data']
		for item in items:
			if item['book_id'] == book_id and item['chapter_start'] == chapter_num:
				try:
					with urllib.request.urlopen(item['path']) as response:
						data = response.read()
						filePath = os.path.join(directory, fileset_id + "." + fileType)
						with open(filePath, "wb") as file:
							file.write(data)
				except urllib.error.URLError as e:
					print("Error downloading the file:", fileset_id, e)
					sys.exit(1)	
	


	def jsonRequest(self, fileset_id, url):
		try:
			with urllib.request.urlopen(url) as response:
				data = response.read()
				try:
					content = json.loads(data.decode('utf-8'))
					return content
				except json.JSONDecodeError:
					print("The file is not json", fileset_id)
					sys.exit(1)
		except urllib.error.URLError as e:
			print("Error downloading the file:", fileset_id, e)
			sys.exit(1)			


bookSeqMap = {'GEN': 1,
'EXO': 2,
'LEV': 3,
'NUM': 4,
'DEU': 5,
'JOS': 6,
'JDG': 7,
'RUT': 8,
'1SA': 9,
'2SA': 10,
'1KI': 11,
'2KI': 12,
'1CH': 13,
'2CH': 14,
'EZR': 15,
'NEH': 16,
'EST': 17,
'JOB': 18,
'PSA': 19,
'PRO': 20,
'ECC': 21,
'SNG': 22,
'ISA': 23,
'JER': 24,
'LAM': 25,
'EZK': 26,
'DAN': 27,
'HOS': 28,
'JOL': 29,
'AMO': 30,
'OBA': 31,
'JON': 32,
'MIC': 33,
'NAM': 34,
'HAB': 35,
'ZEP': 36,
'HAG': 37,
'ZEC': 38,
'MAL': 39,
'MAT': 41,
'MRK': 42,
'LUK': 43,
'JHN': 44,
'ACT': 45,
'ROM': 46,
'1CO': 47,
'2CO': 48,
'GAL': 49,
'EPH': 50,
'PHP': 51,
'COL': 52,
'1TH': 53,
'2TH': 54,
'1TI': 55,
'2TI': 56,
'TIT': 57,
'PHM': 58,
'HEB': 59,
'JAS': 60,
'1PE': 61,
'2PE': 62,
'1JN': 63,
'2JN': 64,
'3JN': 65,
'JUD': 66,
'REV': 67,
'TOB': 68,
'JDT': 69,
'ESG': 70,
'WIS': 71,
'SIR': 72,
'BAR': 73,
'LJE': 74,
'S3Y': 75,
'SUS': 76,
'BEL': 77,
'1MA': 78,
'2MA': 79,
'3MA': 80,
'4MA': 81,
'1ES': 82,
'2ES': 83,
'MAN': 84,
'PS2': 85,
'ODA': 86,
'PSS': 87,
'EZA': 88,
'5EZ': 89,
'6EZ': 90,
'DAG': 91,
'PS3': 92,
'2BA': 93,
'LBA': 94,
'JUB': 95,
'ENO': 96,
'1MQ': 97,
'2MQ': 98,
'3MQ': 100,
'REP': 101,
'4BA': 102,
'LAO': 103,
'FRT': 104,
'BAK': 105,
'OTH': 106,
'INT': 107,
'CNC': 108,
'GLO': 109,
'TDX': 110,
'NDX': 111,
'XXA': 112,
'XXB': 113,
'XXC': 114,
'XXD': 115,
'XXE': 116,
'XXF': 117,
'XXG': 118 }


if __name__ == "__main__":
	dbp = DBP_Adapter()
	dbp.process()
