import os
from transformers import Wav2Vec2ForCTC, AutoProcessor
from datasets import Dataset, Audio
from transformers import Wav2Vec2ForCTC, AutoProcessor
import torch

# Documentation used to write this program
# https://huggingface.co/docs/transformers/main/en/model_doc/mms

class MMSAutoSpeechRecognition:

    def __init__(self):
        #model_id = "facebook/mms-1b-all"
        #target_lang = "npi"
        #self.processor = AutoProcessor.from_pretrained(model_id, target_lang=target_lang)
        #self.model = Wav2Vec2ForCTC.from_pretrained(model_id, target_lang=target_lang, ignore_mismatched_sizes=True)
        model_id = "facebook/mms-1b-all"
        self.processor = AutoProcessor.from_pretrained(model_id)
        self.model = Wav2Vec2ForCTC.from_pretrained(model_id)


    def recognize(self, lang: str, audioFile: str):
        fromDict = Dataset.from_dict({"audio": [audioFile]})
        streamData = fromDict.cast_column("audio", Audio(sampling_rate=16000))
        sample = next(iter(streamData))["audio"]["array"]

        self.processor.tokenizer.set_target_lang(lang)
        self.model.load_adapter(lang)
        inputs = self.processor(sample, sampling_rate=16_000, return_tensors="pt")
        with torch.no_grad():
            outputs = self.model(**inputs).logits
        ids = torch.argmax(outputs, dim=-1)[0]
        transcription = self.processor.decode(ids)
        return transcription


if __name__ == "__main__":
    asr = MMSAutoSpeechRecognition()
    audioFile = os.environ.get("FCBH_DATASET_FILES") + "/NPIDPI/NPIDPIN1DA/B02___01_Mark________NPIDPIN1DA.wav"
    transcription = asr.recognize("npi", audioFile)
    print(transcription)