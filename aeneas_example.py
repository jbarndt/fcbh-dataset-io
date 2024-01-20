# requirements 
# install ffmpeg
# install espeak
# pip install numpy
# pip install aeneas

import os
import subprocess

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


def storeAeneas(audioFile, outputFile):
    resultSet = db.selectAudioFile(audioFile)
    with open(outputFile, 'r') as file:
        timestamps = json.load(file)
        segments = timestamps['fragments']
        if len(segments) != numTextWords:
            print("ERROR: Num Text Words =", len(resultSet), "Num Audio Words =",
                len(segments))
        for index, seg in enumerate(segments):
            (id, word, src_word) = resultSet[index]
            if len(seg['children']) > 0:
                print("Error in segments there are children", seg)
            if len(seg['lines']) != 1:
                print("Error lines is not 1 word", seg)
            elif word != seg['lines'][0]:
                print("Error parsed word and aeneas do not match")
            db.updateMFCC(id, float(seg['begin'], float(seg['end'])
    

# Test1
dir = "../Sandeep_sample1/"
audioFile = os.path.join(dir, "audio.mp3")
textFile = os.path.join(dir, "mplain.txt")
outFile = os.path.join(dir, "poem_single.json")
aeneas("eng", audioFile, textFile, outFile)

# Test2
audioFile = os.path.join(dir, "audio.mp3")
textFile = os.path.join(dir, "words.txt")
outFile = os.path.join(dir, "words_single.json")
aeneas("eng", audioFile, textFile, outFile)