import librosa
import numpy as np
import json
from DBAdapter import *

# Default hop_length = 512
# number_of_frames = number_of_samples / hop_length
# frame_rate = sample_rate / hop_length

def genMFCC(db, audio_file):
	audioData, sampleRate = librosa.load(audio_file)
	print("sampleRate", sampleRate)
	mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
	print("mfccs shape", mfccs.shape)
	# Load your timestamps from the JSON file
	hopLength = 512 # librosa default
	frameRate = sampleRate / hopLength
	resultSet = db.selectTimestamps(audio_file)
	for (id, word, audio_begin_ts, audio_end_ts) in resultSet:
		print(id, word, audio_begin_ts, audio_end_ts)
		startIndex = int(audio_begin_ts * frameRate)
		endIndex = int(audio_end_ts * frameRate)
		# Slice the MFCC data
		segment = mfccs[:, startIndex:endIndex]
		print("start", startIndex, "end", endIndex, "shape", segment.shape)
		#segments.append(segment)
		db.updateMFCC(id, segment)

#def normalize_mfcc(mfccs):
#	mfccs = librosa.util.normalize(mfccs)
#	return mfccs

#def plot_mfcc(mfccs):
#	import matplotlib.pyplot as plt
#	plt.figure(figsize=(10, 4))
#	librosa.display.specshow(mfccs, x_axis='time')
#	plt.colorbar()
#	plt.title('MFCC')
#	plt.tight_layout()
#	plt.show()
#

#gen_segments(1, "../Sandeep_sample1/words_single.json")
#exit(0)

# Test1
db = DBAdapter("ENG", 1, "Sonnet")
mfccs = genMFCC(db, "../Sandeep_sample1/audio.mp3")
