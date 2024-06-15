from transformers import Wav2Vec2ForCTC, AutoProcessor

model_id = "models/mms1b_all.pt"
target_lang = "eng"

processor = AutoProcessor.from_pretrained(model_id, target_lang=target_lang)
model = Wav2Vec2ForCTC.from_pretrained(model_id, target_lang=target_lang, ignore_mismatched_sizes=True)


