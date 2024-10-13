import os
import sys
import re
import json
import torch
import torchaudio
import uroman as ur
from typing import List
from torchaudio.pipelines import MMS_FA as bundle
import ffmpeg

# mms_forced_aligner was developed with the following documentation
# https://pytorch.org/audio/main/tutorials/forced_alignment_for_multilingual_data_tutorial.html
# https://pytorch.org/audio/main/tutorials/ctc_forced_alignment_api_tutorial.html
# A conda environment was created
# conda create -n mms_fa python=3.11
# conda activate mms_fa
# cd mms_fa
# conda install pytorch::pytorch torchvision torchaudio -c pytorch
# pip install uroman # conda does not have it
# pip install soundfile # needed for loading audio files
# pip install ffmpeg-python # needed for audio file conversion to wav


class MMSForcedAligner:

    ## create an instance of a forced aligner to be used any number of times
    def __init__(self):
        if torch.cuda.is_available():
            self.device = "cuda"
        elif torch.backends.mps.is_available():
            self.device = "cpu" ## mps is not yet supported
        else:
            self.device = "cpu"
        print("device", self.device)
        self.model = bundle.get_model(with_star=False)
        self.model.to(self.device)
        self.tokenizer = bundle.get_tokenizer()
        self.aligner = bundle.get_aligner()
        self.uroman = ur.Uroman() # load uroman
        # print(bundle.get_dict())


    ## convert text output from the dataset program to a map keyed on book:chapter
    def loadText(self, textPath:str):
        with open(textPath, 'r') as f:
            text = json.load(f)
        result = {}
        for verse in text:
            ref = verse['book_id'] + ':' + str(verse['chapter_num'])
            list = result.get(ref, [])
            list.append(verse)
            result[ref] = list
        jsonMap = {}
        for ref, verses in result.items():
            jsonMap[ref] = json.dumps(verses, indent=2)
        return jsonMap

    ## convert the text of a chapter to a list of normalized and uroman verse text and references
    def prepareText(self, lang:str, jsonText:str):
        verses = json.loads(jsonText)
        textList = []
        refList = []
        for verse in verses:
            text = verse["script_text"]
            text = self.uroman.romanize_string(text, lcode=lang)
            text = text.lower()
            text = text.replace("’", "'")
            text = re.sub("([^a-z' ])", " ", text)
            text = re.sub(' +', ' ', text)
            text = text.strip()
            words = text.split()
            for i in range(0, len(words), 1):
                textList.append(words[i])
                refList.append(verse["verse_str"] + "\t" + str(i))
        return textList, refList

    ## load a chapter audio file, converting to .wav if needed
    def prepareAudio(self, audioPath: str):
        filename, ext = os.path.splitext(audioPath)
        print("filename", filename, "ext", ext)
        if ext == ".mp3":
            outputFile = filename + ".wav"
            stream = ffmpeg.input(audioPath)
            stream = ffmpeg.output(stream, outputFile, acodec="pcm_s16le", ar=16000)
            stream = ffmpeg.overwrite_output(stream)
            ffmpeg.run(
                stream,
                overwrite_output=True,
                cmd=["ffmpeg", "-loglevel", "error"],  # type: ignore
            )
        elif ext == ".wav":
            outputFile = audioFile
        else:
            print("This audio format is not supported.", audioPath)
            os.exit(1)
        waveform, sample_rate = torchaudio.load(outputFile)#, frame_offset=int(0.5 * bundle.sample_rate), num_frames=int(2.5 * bundle.sample_rate))
        assert sample_rate == bundle.sample_rate
        return waveform, sample_rate

    ## execute force alignment on one Bible chapter
    def align(self, lang: str, book: str, chapter: int, audioPath: str, jsonText: str):
        transcript, references = self.prepareText(lang, jsonText)
        waveform, sample_rate = self.prepareAudio(audioPath)
        with torch.inference_mode():
            emission, _ = self.model(waveform.to(self.device))
            token_spans = self.aligner(emission[0], self.tokenizer(transcript))
        num_frames = emission.size(1)
        ratio = waveform.size(1) / num_frames / sample_rate
        result = []
        assert len(token_spans) == len(transcript)
        assert len(token_spans) == len(references)
        for spans, chars, ref in zip(token_spans, transcript, references):
            timestamp = {}
            timestamp["book"] = book
            timestamp["chapter"] = chapter
            verse, seq = ref.split("\t")
            timestamp["verse"] = verse
            timestamp["seq"] = seq
            timestamp["start"] = round(spans[0].start * ratio, 3)
            timestamp["end"] = round(spans[-1].end * ratio, 3)
            score = sum(s.score * len(s) for s in spans) / sum(len(s) for s in spans)
            timestamp["score"] = round(score, 2)
            timestamp["text"] = chars
            result.append(timestamp)
            print("spans", spans)
            print("timestamp", timestamp)
        return result


    # compute the timestamp and word alignment of each word
    def word_align(self, lang: str, book: str, chapter: int, audioPath: str, jsonText: str):
        result = self.align(lang, book, chapter, audioPath, jsonText)
        return json.dumps(result, indent=2)


    # compute the timestamp and word alignment of each verse
    def verse_align(self, lang: str, book: str, chapter: int, audioPath: str, jsonText: str):
        words = self.align(lang, book, chapter, audioPath, jsonText)
        scores = []
        result = []
        for i in range(0, len(words), 1):
            word = words[i]
            if word["seq"] == "0":
                if i > 0:
                    timestamp["score"] = round(sum(s for s in scores) / len(scores), 2)
                    result.append(timestamp)
                timestamp = word
                word.pop("seq")
                scores.append(word["score"])
            else:
                timestamp["end"] = word["end"]
                scores.append(word["score"])
                timestamp["text"] += " " + word["text"]
        timestamp["score"] = round(sum(s for s in scores) / len(scores), 2)
        result.append(timestamp)
        return json.dumps(result, indent=2)



if __name__ == "__main__":
    print("in main")
    mms_fa = MMSForcedAligner()
    #result = mms_fa.align("deu", "german.flac", "aber seit ich bei ihnen das brot hole")
    audioDir = os.environ.get('FCBH_DATASET_FILES') + "/NPIDPI/NPIDPIN1DA"
    textDir = os.environ.get('FCBH_DATASET_FILES') + "/NPIDPI/NPIDPIN_ET-usx"
    audioFile = audioDir = audioDir + "/B02___01_Mark________NPIDPIN1DA.mp3"
    textFile = "mms_npi.json"
    #waveform, sample_rate = mms_fa.prepareAudio(audioFile)
    #print("len", len(waveform), "rate", sample_rate)
    textMap = mms_fa.loadText(textFile)
    result = mms_fa.verse_align('npi', 'MRK', 1, audioFile, textMap['MRK:1'])
    print(result)


