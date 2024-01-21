# install ffmpeg
# install espeak
# pip install numpy
# pip install aeneas

import os
import json
import subprocess
from DBAdapter import *

# Create a text file of the words using the parsed script in the Database
def createWordsFile(db, audioFile, outputFile):
    with open(outputFile, 'w') as file:
        resultSet = db.selectWordsForFile(audioFile)
        for (id, word, src_word) in resultSet:
            file.write(word + '\n')

# Use Aeneas to produce timestamps for the beginning and ending of each word
def aeneas(language, audioFile, textFile, outputFile):
    command = [
        "python3", "-m", "aeneas.tools.execute_task",
        audioFile,
        textFile,
        f"task_language={language}|os_task_file_format=json|is_text_type=plain",
        outputFile,
        #"-example-words-multilevel --presets-word"
        "-example-words --presets-word"
    ]
    subprocess.run(command)

# Check that the generated output is consistent with the input,
# and store timestamps.
def storeAeneas(db, audioFile, outputFile):
    resultSet = db.selectWordsForFile(audioFile)
    with open(outputFile, 'r') as file:
        timestamps = json.load(file)
        segments = timestamps['fragments']
        if len(segments) != len(resultSet):
            print("ERROR: Num Text Words =", len(resultSet), 
                "Num Audio Words =", len(segments))
        for index, seg in enumerate(segments):
            (id, word, src_word) = resultSet[index]
            if len(seg['children']) > 0:
                print("Error in segments there are children", seg)
            if len(seg['lines']) != 1:
                print("Error lines is not 1 word", seg)
            elif word != seg['lines'][0]:
                print("Error parsed word and aeneas do not match")
            db.updateTimestamps(id, float(seg['begin']), float(seg['end']))
    


# Test1
dir = "../Sandeep_sample1/"
audioFile = os.path.join(dir, "audio.mp3")
textFile = os.path.join(dir, "mplain.txt")
outFile = os.path.join(dir, "poem_single.json")
#aeneas("eng", audioFile, textFile, outFile)

# Test2
textFile = os.path.join(dir, "words.txt")
outFile = os.path.join(dir, "words_single.json")
#aeneas("eng", audioFile, textFile, outFile)

#Test3
textOutFile = "../Sandeep_sample1/sonnet.txt"
if os.path.exists(textOutFile):
    os.remove(textOutFile)
db = DBAdapter("ENG", 1, "Sonnet")
createWordsFile(db, audioFile, textOutFile)
outFile = os.path.join(dir, "sonnet.json")
aeneas("eng", audioFile, textOutFile, outFile)
storeAeneas(db, audioFile, outFile)


