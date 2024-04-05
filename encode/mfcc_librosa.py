import os
import sys
import json
import librosa

# This program puts a json object on stdout.
# Warning: So do not add print statements.

if len(sys.argv) < 2:
    print("Usage: python3 mfcc_librosa.py  audio_file_path", file=sys.stderr)
    sys.exit(1)
audioPath = sys.argv[1]
audioData, sampleRate = librosa.load(audioPath)
mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
hopLength = 512 # librosa default
frameRate = sampleRate / hopLength
result = {}
result["input_file"] = os.path.basename(audioPath)
result["sample_rate"] = sampleRate
result["hop_length"] = hopLength
result["frame_rate"] = frameRate
result["mfcc_shape"] = mfccs.shape
result["mfcc_type"] = str(type(mfccs.dtype))
result["mfccs"] = mfccs.tolist()
print(json.dumps(result))

# python3 mfcc_librosa.py $HOME/FCBH2024/download/ENGWEB/ENGWEBN2DA/B26___01_Jude________ENGWEBN2DA.mp3