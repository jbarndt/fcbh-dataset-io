cd $HOME/FCBH2024/Meta_MMS
git clone https://github.com/pytorch/fairseq
wget -P ./models_new 'https://dl.fbaipublicfiles.com/mms/asr/mms1b_all.pt'

ffmpeg -y -i ../download/ENGWEB/ENGWEBN2DA-mp3-64/B23___01_1John_______ENGWEBN2DA.mp3 -ar 16000 ./audio.wav



export TMPDIR=/Users/gary/FCBH2024/Meta_MMS
export PYTHONPATH=/usr/bin/python3
export PREFIX=INFER
export HYDRA_FULL_ERROR=1
export USER=micro

/usr/bin/python3 fairseq/examples/mms/asr/infer/mms_infer.py --model "mms1b_fl102.pt" --lang "eng" --audio "./audio.wav"

