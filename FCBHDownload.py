import os
import sys
import urllib.request
import json

HOST = "https://4.dbt.io/api/"

class FCBHDownload:


	def process(self):
		print("FCBHDownload version 0.9")
		(isoCode, directory, isoContent) = self.getLanguage()
		cleanContent = self.pruneIsoContent(isoContent)
		if len(cleanContent) > 0:
			transContent = self.displayLanguage(cleanContent)
			bible = transContent['abbr']
			filesetContent = self.displayFilesets(transContent)
			ftype = filesetContent.get('type')
			if ftype == 'text_plain':
				fileset_id = filesetContent['id']
				url = HOST + "download/" + fileset_id + "?v=4&limit=100000"
				content = self.httpRequest(fileset_id, url)
				#sortedContent = sorted(content, key=lambda x: (bookSeqMap[x['book_id']], x['chapter'], x['verse_start']))
				self.saveFile(directory, bible, fileset_id + ".json", content)
			else:
				cloudContent = self.downloadLocation(filesetContent['id'])
				if cloudContent != None:
					self.downloadFiles(directory, cloudContent)


	def getLanguage(self):
		if len(sys.argv) > 2:
			isoCode = sys.argv[1]
			directory = sys.argv[2]
			url = HOST + "bibles?language_code=" + isoCode + "&page=1&limit=100&v=4"
			content = self.httpRequest(isoCode, url)	
			(isoContent, metaContent) = self.parseJson(content)
		else:
			print("Usage: python3 FCBHDownload.py  iosCode  directory")
			print("Requires environment variable: FCBH_DBP_KEY")
			sys.exit(1)
		return (isoCode, directory, isoContent)


	def pruneIsoContent(self, isoContent):
		results = []
		for trans in isoContent:
			filesets = trans['filesets'].get('dbp-prod')
			if filesets != None:
				newFilesets = []
				for fileset in filesets:
					ftype = fileset.get('type')
					if ftype != 'audio_drama_stream' and ftype != 'audio_stream':
						newFilesets.append(fileset)
				if len(newFilesets) > 0:
					trans['filesets']['dbp-prod'] = newFilesets
					results.append(trans)
		return results


	def displayLanguage(self, isoContent):
		first = isoContent[0]
		print()
		print("{: <4}  {: <40}".format(first.get('iso'), first.get('language')))
		print()
		for index, row in enumerate(isoContent):
			iso = row['iso'] if row['iso'] != None else ''
			language = row['language'] if row['language'] != None else ''
			abbr = row['abbr'] if row['abbr'] != None else ''
			name = row['name'] if row['name'] != None else ''
			print("{: <5} {: <10} {: <40}".format(index + 1, abbr, name))
		print()
		if len(isoContent) == 1:
			return isoContent[0]
		else:
			transIndex = 0
			while transIndex < 1 or transIndex > len(isoContent):
				answer = input("Enter number of translation (rtn to exit): ")
				if len(answer) == 0:
					exit(0)
				transIndex = int(answer) if answer.isdigit() else 0
			return isoContent[transIndex - 1]


	def displayFilesets(self, transContent):
		filesets = transContent['filesets']['dbp-prod']
		for index, row in enumerate(filesets):
			filesetId = row['id']
			typ = row['type'] if row['type'] != None else ''
			size = row['size'] if row['size'] != None else ''
			bitrate = row.get('bitrate') if row.get('bitrate') != None else ''
			ftype = row.get('container') if row.get('container') != None else ''
			if typ == 'text_usx':
				ftype = 'usx'
			elif typ == 'text_json':
				ftype = 'json'
			print("{: <5} {: <20} {: <5} {: <15} {: <10} {: <10}".format(index + 1, 
				filesetId, size, typ, ftype, bitrate))
		print()
		fileIndex = 0
		while fileIndex < 1 or fileIndex > len(filesets):
			answer = input("Enter number of fileset (rtn to exit): ")
			if len(answer) == 0:
				exit(0)
			fileIndex = int(answer) if answer.isdigit() else 0
		return filesets[fileIndex - 1]


	def downloadLocation(self, filesetId):
		#print(filesetContent)
		url = HOST + "download/" + filesetId + "?v=4"
		content = self.httpRequest(filesetId, url)
		if content == None:
			return None
		else:
			(json, meta) = self.parseJson(content)
			sortedJson = sorted(json, key=lambda x: (bookSeqMap[x['book_id']], x['chapter_start'], x['verse_start']))
			return sortedJson


	def downloadFiles(self, directory, cloudContent):
		for file in cloudContent:
			url = file.get('path')
			size = file.get('filesize_in_bytes')
			parsedURL = urllib.parse.urlparse(url)
			parts = parsedURL.path.split('/')
			filepath = os.sep.join(parts[2:])
			dirpath = os.path.dirname(filepath)
			filename = os.path.basename(filepath)
			print("Downloading", filepath)
			try:
				with urllib.request.urlopen(url) as response:
					content = response.read()
			except urllib.error.HTTPError as e:
				print("HTTP Error downloading", e.code, e.reason)
			except urllib.error.URLError as e:
				print("Error downloading the file:", param, e)
				sys.exit(1)	
			if content != None:
				if len(content) != size:
					print("Warning for", filepath, "has an expected size of", size, "but, actual size is", len(content))
				self.saveFile(directory, dirpath, filename, content)


	def httpRequest(self, param, url):
		url += "&key=" + os.environ["FCBH_DBP_KEY"]
		try:
			with urllib.request.urlopen(url) as response:
				data = response.read()
				return data
		except urllib.error.HTTPError as e:
			print("HTTP Error downloading", e.code, e.reason)
		except urllib.error.URLError as e:
			print("Error downloading the file:", param, e)
			sys.exit(1)	


	def parseJson(self, content):
		try:
			content = json.loads(content.decode('utf-8'))
			#print("META:", content.get('meta'))
			return (content['data'], content['meta'])
		except json.JSONDecodeError:
			print("The file is not json", param)
			sys.exit(1)	


	def saveFile(self, directory, dirPath, filename, content):
		if content != None:
			fullDir = os.path.join(directory, dirPath)
			if not os.path.exists(fullDir):
				os.makedirs(fullDir)
			filepath = os.path.join(fullDir, filename)
			with open(filepath, "wb") as file:
				file.write(content)	


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
	dbp = FCBHDownload()
	dbp.process()



