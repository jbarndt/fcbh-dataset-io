import torch
import torchaudio
import uroman as ur
import re
from typing import List
from torchaudio.pipelines import MMS_FA as bundle
import time

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
		print("dictionary", bundle.get_dict())
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
		print("sample_rate", sample_rate)
		assert sample_rate == bundle.sample_rate
		return waveform, sample_rate


	def previewWord(self, waveform, spans, num_frames, transcript, sample_rate=bundle.sample_rate):
		ratio = waveform.size(1) / num_frames
		x0 = int(ratio * spans[0].start)
		x1 = int(ratio * spans[-1].end)
		score = sum(s.score * len(s) for s in spans) / sum(len(s) for s in spans)
		print(f"{transcript} ({score:.2f}): {x0 / sample_rate:.3f} - {x1 / sample_rate:.3f} sec")


	def align(self, lang: str, audioPath: str, text: str):
		transcript = self.prepareText(lang, text)
		print("normalized", transcript)
		waveform, sample_rate = self.prepareAudio(audioPath)
		tokens = self.tokenizer(transcript)
		print("tokens", tokens)
		#emission, token_spans = self.compute_alignments(waveform, transcript)
		with torch.inference_mode():
			emission, _ = self.model(waveform.to(self.device))
			token_spans = self.aligner(emission[0], self.tokenizer(transcript))
		#print("emission", emission)
		#print("token_spans", token_spans)
		num_frames = emission.size(1)
		#print("num_frames", num_frames)
		ratio = waveform.size(1) / num_frames / sample_rate
		print("ratio", ratio)
		for spans, chars in zip(token_spans, transcript):
			print("spans", spans)
			t0, t1 = spans[0].start, spans[-1].end
			print("chars", chars, "t0", t0, "t1", t1)
			print("t0ratio", t0 * ratio, "t1*ratio", t1 * ratio)
			score = sum(s.score * len(s) for s in spans) / sum(len(s) for s in spans)
			print("score", score)
			self.previewWord(waveform, spans, num_frames, chars, sample_rate)
			print()



if __name__ == "__main__":
	print("in main")
	mms_fa = MMSForcedAligner()
	mms_fa.align("deu", "german.flac", "aber seit ich bei ihnen das brot hole")

