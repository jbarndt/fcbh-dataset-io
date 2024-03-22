import librosa
import numpy as np
import json
from DBAdapter import *

# Default hop_length = 512
# number_of_frames = number_of_samples / hop_length
# frame_rate = sample_rate / hop_length

def genMFCC(db, audioDir, audio_file):
	audioPath = os.path.join(audioDir, audio_file)
	audioData, sampleRate = librosa.load(audioPath)
	print("sampleRate", sampleRate)
	mfccs = librosa.feature.mfcc(y=audioData, sr=sampleRate, n_mfcc=13)
	#print("mfccs shape", mfccs.shape, "type", type(mfccs.dtype))
	hopLength = 512 # librosa default
	frameRate = sampleRate / hopLength
	resultSet = db.selectWordTimestampsByFile(audio_file)
	#print("words", len(resultSet))
	for (word_id, word, word_begin_ts, word_end_ts) in resultSet:
		print(word_id, word, word_begin_ts, word_end_ts)
		startIndex = int(word_begin_ts * frameRate)
		endIndex = int(word_end_ts * frameRate)
		# Slice the MFCC data
		segment = mfccs[:, startIndex:endIndex]
		print("start", startIndex, "end", endIndex, "shape", segment.shape)
		db.addWordMFCC(word_id, segment)
	db.updateWordMFCCs()


def normPadMFCC(db, normalize):
	mfccTuples = db.selectWordMFCCs() # selects id, MFCCs in numpy
	#print("mfcc", len(mfccTuples))
	mfccList = []
	for (word_id, mfcc) in mfccTuples:
		mfccList.append(mfcc)
	joinedMFCCs = np.concatenate(mfccList, axis=1)
	mean = np.mean(joinedMFCCs, axis=1)
	stds = np.std(joinedMFCCs, axis=1)
	maxLen = max(array.shape[1] for array in mfccList)
	for (word_id, mfcc) in mfccTuples:
		if normalize:
			mfcc2 = (mfcc - mean[:, None]) / stds[:, None]
		else:
			mfcc2 = mfcc
		padded = np.pad(mfcc, ((0, 0), (0, maxLen - mfcc2.shape[1])), 'constant')
		db.addPadWordMFCC(word_id, padded)
	db.updatePadWordMFCCs()


'''
# Test2
db = DBAdapter("ENG", 3, "Excel")
dir = "../../Desktop/Mark_Scott_1_1-31-2024/Audio Files"
audioFile = "N2_MZI_BSM_046_LUK_002_VOX.wav"
#audioPath = os.path.join(dir, audioFile)
genMFCC(db, dir, audioFile)
normPadMFCC(db, True)
'''

