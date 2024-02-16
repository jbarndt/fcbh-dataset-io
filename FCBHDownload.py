import os
import sys
import urllib.request
import json

HOST = "https://4.dbt.io/api/"
#CONFIG = os.path.join(os.environ["HOME"], "FCBHDownload.cfg")

class FCBHDownload:


	def process(self):
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
				#dirPath = os.path.join(directory, bible)
				#os.path.makedirs(dirpath)
				self.saveFile(directory, bible, fileset_id + ".json", content)
			else:
				cloudContent = self.downloadLocation(filesetContent['id'])
				if cloudContent != None:
					self.downloadFiles(directory, cloudContent)
				#print(cloudContent)


	def getLanguage(self):
		if len(sys.argv) > 2:
			isoCode = sys.argv[1]
			directory = sys.argv[2]
			url = HOST + "bibles?language_code=" + isoCode + "&page=1&limit=100&v=4"
			content = self.httpRequest(isoCode, url)	
			isoContent = self.parseJson(content)
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
		#iso = first.get('iso')
		#language = first.get('language')
		print()
		print("{: <4}  {: <40}".format(first.get('iso'), first.get('language')))
		print()
		for index, row in enumerate(isoContent):
			#print(row)
			iso = row['iso'] if row['iso'] != None else ''
			language = row['language'] if row['language'] != None else ''
			abbr = row['abbr'] if row['abbr'] != None else ''
			name = row['name'] if row['name'] != None else ''
			print("{: <5} {: <10} {: <40}".format(index + 1, abbr, name))
		#transIndex = input("Enter number of translation:")
		print()
		if len(isoContent) == 1:
			return isoContent[0]
		else:
			transIndex = 0
			while transIndex < 1 or transIndex > len(isoContent):
				answer = input("Enter number of translation: ")
				transIndex = int(answer) if answer.isdigit() else 0
			#return transIndex - 1
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
			answer = input("Enter number of fileset: ")
			fileIndex = int(answer) if answer.isdigit() else 0
		return filesets[fileIndex - 1]


#	def prepareDirectory(self, directory):
#		directory = input("Enter directory to store fileset: ")
#		if not os.path.exists(directory):
#			os.mkdir(directory)
#		return directory


	def downloadLocation(self, filesetId):
		#print(filesetContent)
		url = HOST + "download/" + filesetId + "?v=4"
		content = self.httpRequest(filesetId, url)
		if content == None:
			return None
		else:
			json = self.parseJson(content)
			return json


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
			#content = self.httpRequest(filepath, url)
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


#	def readISOFile(self):
#		with open("FCBHDownload.cfg", "r") as file:
#			content = file.read()
#			return content


	def parseJson(self, content):
		try:
			content = json.loads(content.decode('utf-8'))
			#print("META:", content.get('meta'))
			return content['data']
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


if __name__ == "__main__":
	dbp = FCBHDownload()
	dbp.process()
