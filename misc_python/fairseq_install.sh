# I only did

/usr/bin/pip3 install fairseq


#cd $HOME/FCBH2024/Meta_MMS

#git clone https://github.com/pytorch/fairseq
#cd fairseq
#CFLAGS="-stdlib=libc++" /usr/bin/pip3 install --editable ./

# to install the latest stable release (0.10.x)
# pip install fairseq


#PATH changed to include /

#/usr/bin/pip3 install --upgrade pip
#/usr/bin/pip3 install --upgrade setuptools wheel

cd $HOME/FCBH2024/Meta_MMS

git clone https://github.com/pytorch/fairseq

cd fairseq

export TMPDIR=/Users/gary/tmp_dir
export PYTHONPATH=/usr/bin/python3
export PREFIX=INFER
export HYDRA_FULL_ERROR=1
export USER=micro

/usr/bin/python3 examples/mms/asr/infer/mms_infer.py --model "../models/mms1b_all.pt" --lang "eng" \
  --audio "../audio.wav"
