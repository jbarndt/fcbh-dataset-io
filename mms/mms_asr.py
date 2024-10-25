import os
import sys
import torch
from transformers import Wav2Vec2ForCTC, AutoProcessor
from datasets import Dataset, Audio

## Documentation used to write this program
## https://huggingface.co/docs/transformers/main/en/model_doc/mms

def isSupportedLanguage(modelId:str, lang:str):
    processor = AutoProcessor.from_pretrained(modelId)
    dict = processor.tokenizer.vocab.keys()
    for l in dict:
        if l == lang:
            return True
    return False


if len(sys.argv) < 2:
    print("Usage: mms_asr.py  {iso639-3}")
    sys.exit(1)
lang = sys.argv[1]
modelId = "facebook/mms-1b-all"
if not isSupportedLanguage(modelId, lang):
    print(lang, "is not supported by", modelId)
processor = AutoProcessor.from_pretrained(modelId, target_lang=lang)
model = Wav2Vec2ForCTC.from_pretrained(modelId, target_lang=lang, ignore_mismatched_sizes=True)

for line in sys.stdin:
    audioFile = line.strip()
    fromDict = Dataset.from_dict({"audio": [audioFile]})
    streamData = fromDict.cast_column("audio", Audio(sampling_rate=16000))
    sample = next(iter(streamData))["audio"]["array"]

    inputs = processor(sample, sampling_rate=16_000, return_tensors="pt")
    with torch.no_grad():
        outputs = model(**inputs).logits
    ids = torch.argmax(outputs, dim=-1)[0]
    transcription = processor.decode(ids)
    sys.stdout.write(transcription)
    sys.stdout.write("\n")
    sys.stdout.flush()


#$HOME + "/NPIDPI/NPIDPIN1DA/B02___01_Mark________NPIDPIN1DA.wav"

# /Users/gary/FCBH2024/download/NPIDPI/NPIDPIN1DA/B02___01_Mark________NPIDPIN1DA.wav

## Testing
## conda activate mms_fa
## python mms_asr.py  npi
## /Users/gary/FCBH2024/download/NPIDPI/NPIDPIN1DA/B02___01_Mark________NPIDPIN1DA.wav
## ctrl-D

