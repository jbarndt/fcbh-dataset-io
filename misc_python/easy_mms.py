

# https://abdeladim-s.github.io/easymms/#easymms.models.alignment.AlignmentModel

# https://github.com/abdeladim-s/easymms/blob/main/README.md


# /usr/bin/pip3 install easymms

# /usr/bin/pip3 install -U --pre torchaudio --index-url https://download.pytorch.org/whl/nightly/cu118

# /usr/bin/pip3 install tensorboardX
import os
from easymms.models.asr import ASRModel
import locale

locale.getpreferredencoding = lambda: "UTF-8"

data = os.environ['FCBH_DATASET_DB']

#model = data + '/easy_mms/models/mms1b_all.pt'
model = data + '/easy_mms/models/mms1b_fl102.pt'
print("model", model)
file1 = data + '/download/ENGWEB/ENGWEBN2DA-mp3-64/B25___01_3John_______ENGWEBN2DA.mp3'
print("file", file1)
asr = ASRModel(model=model)
files = [file1]
transcriptions = asr.transcribe(files, lang='eng', align=False)
for i, transcription in enumerate(transcriptions):
    #print(f">>> file {files[i]}")
    print(transcription)


# /usr/bin/pip3 install --upgrade fairseq

# /usr/bin/pip3 install --upgrade pytorch

# cd /path/to/fairseq-py/
# python examples/mms/asr/infer/mms_infer.py --model "/path/to/asr/model" --lang lang_code \
#  --audio "/path/to/audio_1.wav" "/path/to/audio_2.wav" "/path/to/audio_3.wav"