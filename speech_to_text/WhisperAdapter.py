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

'''
1. Retry Whisper with a non-drama version of the audio.
2. Pi recommends preprocessing in audacity
3. Do not get text, but get segments, and iterate over it.
4. Capture start, end timestamps
5. Capture text
6. Capture tokens, if I have a place for it
7. Capture avg_logprob
8. Capture no_speech_prob
9. Capture compression ratio

'segments': [{'id': 0, 'seek': 0, 'start': 0.0, 'end': 3.24, 'text': ' Chapter 3', 'tokens': [50363, 7006, 513, 50525],
'temperature': 0.0, 'avg_logprob': -0.2316610102067914, 'compression_ratio': 1.46, 'no_speech_prob': 
0.2119932472705841}, 

'''

class WhisperAdapter:


	def __init__(self, db):
		self.db = db
		self.model = whisper.load_model("medium.en") # "small" "base" "medium" "large" are options
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
					#if file.startswith("B17"): ## TITUS
					if file.startswith("B"): ## New Testament
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
	if len(sys.argv) < 3:
		print("Usage: WhisperAdapter  bibleId  audio_filesetId")
		sys.exit(1)
	bibleId = sys.argv[1]
	filesetId = sys.argv[2]
	db = DBAdapter(bibleId + "_WHISPER.db")
	whisp = WhisperAdapter(db)
	whisp.insertIdent(bibleId)
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





