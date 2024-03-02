import os
import sys
import whisper
from DBAdapter import *
from WordParser import *
from Booknames import *

#https://github.com/openai/whisper

#model = whisper.load_model("base")
#result = model.transcribe("audio.mp3")
#print(result["text"])

class WhisperAdapter:


	def __init__(self, db):
		self.db = db
		self.model = whisper.load_model("medium") # "small" "base" "medium" "large" are options
		self.books = Booknames()


	def insertIdent(self, bibleId):
		iso = bibleId[:3]
		version = bibleId[3:]
		self.db.insertIdent(bibleId, iso, version, "WHISPER", None, None, None, None, None)


	def processDirectory(self, directory):
		for file in sorted(os.listdir(directory)):
			if not file.startswith("."):
				print(file)
				resultSet = self.db.selectScriptsByFile(file)
				if len(resultSet) == 0:
					if file.startswith("B17"): ## TITUS
						self.processFile(directory, file)


	def processFile(self, directory, filename):
		filePath = os.path.join(directory, filename)
		(bookId, chapter) = self.parseFilename(filename)
		result = self.model.transcribe(filePath)
		scriptText = result["text"]
		self.db.addScript(bookId, chapter, filename, 1, None, None, None, None, scriptText)
		self.db.insertScripts()


	def parseFilename(self, filename):
		chapter = int(filename[6:8])
		bookName = filename[9:21].replace("_", "")
		bookId = self.books.usfmBookId(bookName)
		return (bookId, chapter)


if __name__ == "__main__":
	if len(sys.argv) < 2:
		print("Usage: WhisperAdapter  bibleId")
		sys.exit(1)
	bibleId = sys.argv[1]
	db = DBAdapter(bibleId + "_WHISPER.db")
	whisp = WhisperAdapter(db)
	whisp.insertIdent(bibleId)
	filesetId = bibleId + "N2DA"
	directory = os.path.join(os.environ['FCBH_DATASET_FILES'], bibleId, filesetId)
	whisp.processDirectory(directory)
	file = WordParser(db)
	file.parse()
	db.close()

'''
pip3 install git+https://github.com/openai/whisper.git 

WARNING: The script isympy is installed in '/Users/gary/Library/Python/3.9/bin' which is not on PATH.
Consider adding this directory to PATH or, if you prefer to suppress this warning, use --no-warn-script-location.

WARNING: The scripts convert-caffe2-to-onnx, convert-onnx-to-caffe2 and torchrun are installed in '/Users/gary/Library/Python/3.9/bin' which is not on PATH.
Consider adding this directory to PATH or, if you prefer to suppress this warning, use --no-warn-script-location.

WARNING: The script whisper is installed in '/Users/gary/Library/Python/3.9/bin' which is not on PATH.
Consider adding this directory to PATH or, if you prefer to suppress this warning, use --no-warn-script-location.
'''





