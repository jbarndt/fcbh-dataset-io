import os
import sys
import json
import librosa

# This program puts a json object on stdout.
# Warning: So do not add print statements.

if len(sys.argv) < 3:
    print("Usage: python3 mfcc_librosa.py  audio_file_path  n_mfcc", file=sys.stderr)
    print("n_mfcc is the number of features, usually 13 to 20", file=sys.stderr)
    sys.exit(1)
audioPath = sys.argv[1]
nMfcc = int(sys.argv[2])
audioData, sampleRate = librosa.load(audioPath)
mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=nMfcc)
mfccs_T = mfccs.T
hopLength = 512 # librosa default
frameRate = sampleRate / hopLength
result = {}
result["input_file"] = os.path.basename(audioPath)
result["sample_rate"] = sampleRate
result["hop_length"] = hopLength
result["frame_rate"] = frameRate
result["mfcc_shape"] = mfccs_T.shape
result["mfcc_type"] = str(type(mfccs_T.dtype))
result["mfccs"] = mfccs_T.tolist()
print(json.dumps(result))

# python mfcc_librosa.py $HOME/FCBH2024/download/ENGWEB/ENGWEBN2DA-mp3-64/B26___01_Jude________ENGWEBN2DA.mp3 20