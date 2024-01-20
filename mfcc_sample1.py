import librosa
import numpy as np
import json

# Default hop_length = 512
# number_of_frames = number_of_samples / hop_length
# frame_rate = sample_rate / hop_length

def gen_mfcc(audioFile, timestampFile):
	audioData, sampleRate = librosa.load(audioFile)
	print("sampleRate", sampleRate)
	mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
	print("mfccs shape", mfccs.shape)
	# Load your timestamps from the JSON file
	hopLength = 512 # librosa default
	frameRate = sampleRate / hopLength
	with open(timestampFile, 'r') as file:
		timestamps = json.load(file)
		segments = []
		for seg in timestamps['fragments']:
			beginTS = float(seg['begin'])
			endTS = float(seg['end'])
			print(seg['id'], beginTS, endTS, seg['lines'], len(seg['children']))
			if len(seg['children']) > 0: ## move check to generation process
				print("Error in segments there are children", seg)
			startIndex = int(beginTS * frameRate)
			endIndex = int(endTS * frameRate)
			# Slice the MFCC data
			segment = mfccs[:, startIndex:endIndex]
			print("start", startIndex, "end", endIndex, "shape", segment.shape)
			segments.append(segment)
			#mysql.updateMccs(id, beginTS, endTS, segment)

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
mfccs = gen_mfcc("../Sandeep_sample1/audio.mp3", "../Sandeep_sample1/words_single.json")
#mfccs2 = normalize_mfcc(mfccs)
