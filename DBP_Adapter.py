
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
		try:
			with urllib.request.urlopen(url) as response:
				data = response.read()  # Read the response
				try:
					content = json.loads(data.decode('utf-8'))
					print("The file is a JSON. Loaded into memory.")
					return content
				except json.JSONDecodeError:
					print("The file is not json", fileset_id)
					sys.exit(1)
		except urllib.error.URLError as e:
			print("Error downloading the file:", fileset_id, e)
			sys.exit(1)



	def createDatabase(self, fileset_id):
		isoCode = fileset_id[:3]
		versionCode = fileset_id[3:6]
		database = isoCode + "_1_" + versionCode + ".db"
		if os.path.exists(database):
			os.remove(database)
		db = DBAdapter(isoCode, 1, versionCode)
		return db


	def loadAudioScript(self, db, jsonArray):
		script_num = 0
		for rec in jsonArray:
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


	def downloadFile(directory, fileType, fileset_id, book_id, chapter_num):
		a = 2


#	def httpRequest(self, url)
#		try:
#			response = requests.get(url)
			# Raise an exception if the request was unsuccessful
#			response.raise_for_status()
#			try:
#				content = json.loads(response.text)
#				print("The file is a JSON. Loaded into memory.")
#				return content  # Return the JSON content if successful
#			except json.JSONDecodeError:
				# If it's not JSON, save it to a file
#				filename = url.split('/')[-1]  # Extract filename from URL
#				with open(filename, 'wb') as f:
#					f.write(response.content)
#				print(f"The file is not a JSON. Saved as {filename}.")
#
#		except requests.RequestException as e:
#			print(f"Error downloading the file: {e}")


#def download_file(url):
#    try:
#        # Send a GET request to the URL
#        with urllib.request.urlopen(url) as response:
#            data = response.read()  # Read the response
#
#            # Try to decode the response as UTF-8 and load it as JSON
#            try:
#                content = json.loads(data.decode('utf-8'))
#                print("The file is a JSON. Loaded into memory.")
#                return content
#            except json.JSONDecodeError:
#                # If it's not JSON, save it to a file
#                filename = url.split('/')[-1]  # Extract filename from URL
#                with open(filename, 'wb') as f:
#                    f.write(data)
#                print(f"The file is not a JSON. Saved as {filename}.")
#
#    except urllib.error.URLError as e:
#        print(f"Error downloading the file: {e}")
#
# Example usage
#url = 'https://example.com/somefile.json'  # Replace this with the actual URL
#downloaded_content = download_file(url)

# If it's a JSON file, you can now work with `downloaded_content` as a dictionary or list


if __name__ == "__main__":
	dbp = DBP_Adapter()
	dbp.process()
