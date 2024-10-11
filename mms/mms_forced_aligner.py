import torch
import torchaudio
import uroman as ur
import re
from typing import List
from torchaudio.pipelines import MMS_FA as bundle

# mms_forced_aligner was developed with the following documentation
# https://pytorch.org/audio/main/tutorials/forced_alignment_for_multilingual_data_tutorial.html
# A conda environment was created
# conda create -n mms_fa python=3.11
# conda activate mms_fa
# cd mms_fa
# conda install pytorch::pytorch torchvision torchaudio -c pytorch
# pip install uroman # conda does not have it
# pip install soundfile # needed for loading audio files

class MMSForcedAligner:

	def __init__(self):
		if torch.cuda.is_available():
			self.device = "cuda"
		elif torch.backends.mps.is_available():
			self.device = "cpu" ## mps is not yet supported
		else:
			self.device = "cpu"
		print("device", self.device)
		self.model = bundle.get_model()
		self.model.to(self.device)
		self.tokenizer = bundle.get_tokenizer()
		self.aligner = bundle.get_aligner()
		self.uroman = ur.Uroman() # load uroman


	def prepareText(self, lang:str, text:str):
		text = self.uroman.romanize_string(text, lcode=lang)
		text = text.lower()
		text = text.replace("’", "'")
		text = re.sub("([^a-z' ])", " ", text)
		text = re.sub(' +', ' ', text)
		text = text.strip()
		print("text", text)
		return text.split()


	def prepareAudio(self, audioPath: str):
		waveform, sample_rate = torchaudio.load(audioPath, frame_offset=int(0.5 * bundle.sample_rate), num_frames=int(2.5 * bundle.sample_rate))
		assert sample_rate == bundle.sample_rate
		return waveform, sample_rate


	def align(self, lang: str, audioPath: str, text: str):
		transcript = self.prepareText(lang, text)
		waveform, sample_rate = self.prepareAudio(audioPath)
		tokens = self.tokenizer(transcript)
		with torch.inference_mode():
			emission, _ = self.model(waveform.to(self.device))
			token_spans = self.aligner(emission[0], self.tokenizer(transcript))
		num_frames = emission.size(1)
		ratio = waveform.size(1) / num_frames / sample_rate
		result = []
		for spans, chars in zip(token_spans, transcript):
			timestamp = {}
			timestamp["text"] = chars
			timestamp["start"] = round(spans[0].start * ratio, 3)
			timestamp["end"] = round(spans[-1].end * ratio, 3)
			score = sum(s.score * len(s) for s in spans) / sum(len(s) for s in spans)
			timestamp["score"] = round(score, 2)
			result.append(timestamp)
			print("spans", spans)
			print("timestamp", timestamp)
			print()
		return json.dumps(result, indent=2)


if __name__ == "__main__":
	print("in main")
	mms_fa = MMSForcedAligner()
	result = mms_fa.align("deu", "german.flac", "aber seit ich bei ihnen das brot hole")
	print(result)

