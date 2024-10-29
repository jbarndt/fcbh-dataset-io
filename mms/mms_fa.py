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

# convert the text of a chapter to a list of normalized and uroman verse text and references
def prepareText(lang:str, verses):
    textList = []
    refList = []
    for verse in verses:
        text = verse["text"]
        text = uroman.romanize_string(text, lcode=lang)
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

# load a chapter audio file, converting to .wav if needed
def prepareAudio(audioPath: str):
    filename, ext = os.path.splitext(audioPath)
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
        outputFile = audioPath
    else:
        sys.stderr.write("This audio format is not supported. " + audioPath + "\n")
        os.exit(1)
    waveform, sample_rate = torchaudio.load(outputFile)#, frame_offset=int(0.5 * bundle.sample_rate), num_frames=int(2.5 * bundle.sample_rate))
    assert sample_rate == bundle.sample_rate
    return waveform, sample_rate

# execute force alignment on one Bible chapter
def align(book: str, chapter: int, audioPath: str, verses):
    transcript, references = prepareText(lang, verses)
    waveform, sample_rate = prepareAudio(audioPath)
    with torch.inference_mode():
        emission, _ = model(waveform.to(device))
        token_spans = aligner(emission[0], tokenizer(transcript))
    num_frames = emission.size(1)
    ratio = waveform.size(1) / num_frames / sample_rate
    result = []
    assert len(token_spans) == len(transcript)
    assert len(token_spans) == len(references)
    for spans, chars, ref in zip(token_spans, transcript, references):
        timestamp = {}
        #timestamp["book_id"] = book
        #timestamp["chapter_num"] = chapter
        verse, seq = ref.split("\t")
        timestamp["verse_str"] = verse
        timestamp["word_seq"] = int(seq)
        timestamp["begin_ts"] = round(spans[0].start * ratio, 3)
        timestamp["end_ts"] = round(spans[-1].end * ratio, 3)
        score = sum(s.score * len(s) for s in spans) / sum(len(s) for s in spans)
        timestamp["fa_score"] = round(score, 2)
        timestamp["uroman"] = chars
        #timestamp["audio_file"] = os.path.basename(audioPath)
        result.append(timestamp)
    return result

if len(sys.argv) < 2:
    sys.stderr.write("Usage: mms_fa.py  {iso639-3}\n")
    sys.exit(1)
lang = sys.argv[1]
if torch.cuda.is_available():
    device = "cuda"
elif torch.backends.mps.is_available():
    device = "cpu" ## mps is not yet supported
else:
    device = "cpu"
model = bundle.get_model(with_star=False)
model.to(device)
tokenizer = bundle.get_tokenizer()
aligner = bundle.get_aligner()
uroman = ur.Uroman() # load uroman
for line in sys.stdin:
    inp = json.loads(line)
    results = align(inp["book_id"], inp["chapter"], inp["audio_file"], inp["verses"])
    output = json.dumps(results)
    sys.stdout.write(output)
    sys.stdout.write("\n")
    sys.stdout.flush()

# Testing
# conda activate mms_fa
# time python mms_fa.py eng < engweb_fa_inp.json > engweb_fa_out.json


