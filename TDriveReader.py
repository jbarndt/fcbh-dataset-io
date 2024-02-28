import re
import os 
import sys
import urllib.request
import json

'''
This is a one-time program.  
The program reads a directory list of audio data,
It finds the xlsm files that are for OT or NT, and finds the filesetId
It then looks up the bibleId in DBP, and identifies the bibles
that have USX + plain_text + audio
It produces a file containing the ones found.

34561 lines in file
925 xlsm files
616 filesetIds
419 are N or O
285 distinct filesetIds
'''

HOST = "https://4.dbt.io/api/"


def readFileListing(filename):
	result = {}
	with open(filename, "r") as file:
		content = file.read()
	for line in content.split("\n"):
		if line.endswith(".xlsm"):
			#print(line)
			filesetId = parseFilesetId(line)
			if filesetId != None:
				if filesetId[0] == "N" or filesetId[0] == "O":
					filenames = result.get(filesetId, [])
					filenames.append(line)
					result[filesetId] = filenames
	return result


def parseFilesetId(line):
	pattern = re.compile(r"[A-Z][0-9][A-Z]{6}")
	match = pattern.findall(line)
	if len(match) > 0:
		return match[0]
	else:
		return None


def findDBPFilesets(bibleId):
	print("----", bibleId)
	url = HOST + "bibles/" + bibleId + "?&v=4"
	print(url)
	content = httpRequest(bibleId, url)
	if content == None:
		return False
	dataContent = parseJson(content)
	if dataContent == None:
		return False
	isUseful = checkIfUseful(dataContent)
	if isUseful:
		print("*** is GOOD ***")
	else:
		print("NO GOOD")
	return isUseful


def checkIfUseful(content):
	hasUSX = False
	hasText = False
	hasAudio = False
	for item in content:
		fid = item.get('id')
		typ = item.get('type')
		size = item.get('size')
		print(fid, typ, size)
		if typ == 'text_usx' and (size == 'NT' or size == 'C'):
			hasUSX = True
		if typ == 'text_plain' and (size == 'NT' or size == 'C'):
			hasText = True 
		if (typ == 'audio' or typ == 'audio_drama') and (size == 'NT' or size == 'C'):
			hasAudio = True
	return (hasUSX and hasText and hasAudio)


def httpRequest(param, url):
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


def parseJson(content):
	try:
		content = json.loads(content.decode('utf-8'))
		data = content.get('data')
		filesets = data.get('filesets')
		if filesets == None:
			return None
		dbpProd = filesets.get('dbp-prod')
		return dbpProd
	except json.JSONDecodeError:
		print("The file is not json", param)
		sys.exit(1)	



with open("USEFUL_FILESETS.txt", "w") as file:
	foundMap = readFileListing("../tdrive.txt")
	count = 0
	for filesetId in foundMap.keys():
		bibleId = filesetId[2:]
		isGood = findDBPFilesets(bibleId)
		if isGood:
			file.write("------\n")
			file.write(filesetId + "  " + bibleId + "\n")
			lines = foundMap[filesetId]
			for line in lines:
				file.write(line + "\n")
		count += 1
		#if count == 50:
		#	file.close()
		#	sys.exit(0)
	print(count)
	file.close()




