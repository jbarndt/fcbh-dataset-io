
import os
import sys
import urllib.request
import json

HOST = "https://4.dbt.io/api/"
CONFIG = os.path.join(os.environ["HOME"], "FCBHDownload.cfg")

class FCBHDownload:



#curl "https://4.dbt.io/api/bibles?language_code=tha&page=1&limit=100&v=4&key=b4715786-9b8e-4fbe-a9b9-ff448449b81b"

	#def __init__(self):
		#self.isoContent = None

	def process(self):
		isoContent = self.getLanguage()
		isoContent = self.pruneIsoContent(isoContent)
		transIndex = self.displayLanguage(isoContent)
		filesetIndex = self.displayFilesets(transIndex, isoContent)


	def getLanguage(self):
		print(sys.argv)
		if len(sys.argv) > 1:
			isoCode = sys.argv[1]
			url = HOST + "bibles?language_code=" + isoCode + "&page=1&limit=100&v=4&key=" + os.environ["FCBH_DBP_KEY"]
			content = self.jsonRequest(isoCode, url)	
		elif os.path.exists(CONFIG):
			content = self.readISOFile(CONFIG)
		else:
			print("Usage: python3 FCBHDownload.py  iosCode")
			print("Requires environment variable: FCBH_DBP_KEY")
			sys.exit(1)
		isoContent = self.parseJson(content)
		return isoContent


	def pruneIsoContent(self, isoContent):
		results = []
		for trans in isoContent:
			filesets = trans['filesets'].get('dbp-prod')
			##if filesets != None and len(filesets) > 0:
			if filesets != None:
				newFilesets = []
				for fileset in filesets:
					if fileset.get('type') != 'audio_drama_stream':
						newFilesets.append(fileset)
				if len(newFilesets) > 0:
					trans['filesets'].get('dbp-prod') = newFilesets
					results.append(trans)
		return results


	def displayLanguage(self, isoContent):
		for index, row in enumerate(isoContent):
			#print(row)
			iso = row['iso'] if row['iso'] != None else ''
			language = row['language'] if row['language'] != None else ''
			abbr = row['abbr'] if row['abbr'] != None else ''
			name = row['name'] if row['name'] != None else ''
			print("{: <4} {: <4} {: <10} {: <10} {: <40}".format(index + 1, iso, language, abbr, name))
		#transIndex = input("Enter number of translation:")
		transIndex = 0
		while transIndex < 1 or transIndex > len(isoContent):
			answer = input("Enter number of translation: ")
			transIndex = int(answer) if answer.isdigit() else 0
		return transIndex - 1


	def displayFilesets(self, transIndex, isoContent):
		filesets = isoContent[transIndex]['filesets']['dbp-prod']
		for index, row in enumerate(filesets):
			filesetId = row['id']
			typ = row['type'] if row['type'] != None else ''
			size = row['size'] if row['size'] != None else ''
			bitrate = row.get('bitrate') if row.get('bitrate') != None else ''
			container = row.get('container') if row.get('container') != None else ''
			print("{: <4} {: <20} {: <5} {: <15} {: <10} {: <10}".format(index + 1, 
				filesetId, size, typ, container, bitrate))
		fileIndex = 0
		while fileIndex < 1 or fileIndex > len(filesets):
			answer = input("Enter number of fileset: ")
			fileIndex = int(answer) if answer.isdigit() else 0
		return fileIndex - 1


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
	

	def jsonRequest(self, param, url):
		try:
			with urllib.request.urlopen(url) as response:
				data = response.read()
				return data
		except urllib.error.URLError as e:
			print("Error downloading the file:", param, e)
			sys.exit(1)	


	def readISOFile(self):
		with open("FCBHDownload.cfg", "r") as file:
			content = file.read()
			return content


	def parseJson(self, content):
		try:
			content = json.loads(content.decode('utf-8'))
			return content['data']
		except json.JSONDecodeError:
			print("The file is not json", param)
			sys.exit(1)			





if __name__ == "__main__":
	dbp = FCBHDownload()
	dbp.process()
