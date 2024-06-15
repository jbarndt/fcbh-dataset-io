import torch
import torchaudio
from transformers import Wav2Vec2Processor

# Load the model
model_path = "/Users/gary/FCBH2024/Meta_MMS/mms1b_all.pt"
model = torch.load(model_path)
model.eval()

# Define the processor
processor = Wav2Vec2Processor.from_pretrained("facebook/mms-processor")

# Load an example audio file
# Replace 'path_to_your_audio_file.wav' with the path to your audio file
audio_input, sample_rate = torchaudio.load("path_to_your_audio_file.wav")

# Preprocess the audio file
input_values = processor(audio_input, sampling_rate=sample_rate, return_tensors="pt").input_values

# Perform inference
with torch.no_grad():
    logits = model(input_values).logits

# Decode the predicted IDs to text
predicted_ids = torch.argmax(logits, dim=-1)
transcription = processor.batch_decode(predicted_ids)[0]

print("Transcription:", transcription)

#!mkdir "temp_dir"
#!git clone https://github.com/pytorch/fairseq
#!wget -P ./models_new 'https://dl.fbaipublicfiles.com/mms/asr/mms1b_all.pt'
#!ffmpeg -y -i ./audio_samples/audio.mp3 -ar 16000 ./audio_samples/audio.wav
#import os

#os.environ["TMPDIR"] = '/content/temp_dir'
#os.environ["PYTHONPATH"] = "."
#os.environ["PREFIX"] = "INFER"
#os.environ["HYDRA_FULL_ERROR"] = "1"
#os.environ["USER"] = "micro"

#!python examples/mms/asr/infer/mms_infer.py --model "/content/fairseq/models_new/mms1b_fl102.pt" --lang "eng" --audio "/content/fairseq/audio_samples/audio.wav"

