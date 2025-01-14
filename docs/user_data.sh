#!/bin/bash -x

# Switch to ec2-user for the remaining operations
sudo -u ec2-user bash << 'EOF'
export HOME=/home/ec2-user
cd $HOME

/opt/conda/bin/conda init bash
source /home/ec2-user/.bashrc

# ---------- base ----------

conda activate base
# /opt/conda/bin/

# Install executables using conda into the base environment
conda install -y ffmpeg -c conda-forge
conda install -y sox -c conda-forge
conda install -y sqlite -c conda-forge

# ---------- aeneas ----------

# https://github.com/readbeyond/aeneas

conda create -y -n aeneas python=3.8
# /opt/conda/envs/aeneas/bin

conda activate aeneas

conda install -y ffmpeg -c conda-forge  ## appears to be redundant

conda install -y numpy -c conda-forge

sudo yum update -y
sudo yum -y install espeak espeak-devel

conda install -y "setuptools <60"

pip install aeneas

python -m aeneas.diagnostics

conda deactivate

# ---------- easy_mms ----------

# conda create -y -n easy_mms python=3.8

# conda activate easy_mms

# pip install easymms

# conda deactivate

# ---------- fasttext ----------

conda create -y -n fasttext --no-default-packages

conda activate fasttext

cd /opt/conda/envs/fasttext
# on Mac
# cd $HOME/miniforge3/envs/fasttext/

git clone https://github.com/facebookresearch/fastText.git
cd fastText
make

conda deactivate

# ---------- librosa ----------

conda create -y -n librosa python=3.11

conda activate librosa

pip install librosa

conda deactivate

# ---------- mms_asr ----------

conda create -y -n mms_asr python=3.11

conda activate mms_asr

conda install -y pytorch torchaudio pytorch-cuda=12.1 -c pytorch -c nvidia
# On Mac
# conda install -y pytorch::pytorch torchaudio -c pytorch

pip install accelerate
pip install datasets
pip install --upgrade transformers
pip install soundfile
pip install librosa

pip install uroman
cp /opt/conda/envs/mms_asr/bin/uroman /opt/conda/envs/mms_asr/bin/uroman.pl
# on Mac
# cp /Users/gary/miniforge3/envs/mms_asr/bin/uroman /Users/gary/miniforge3/envs/mms_asr/bin/uroman.pl

conda deactivate

# ---------- mms_fa ----------

# https://pytorch.org/audio/main/tutorials/forced_alignment_for_multilingual_data_tutorial.html

conda create -y -n mms_fa python=3.11

conda activate mms_fa

conda install -y pytorch torchaudio pytorch-cuda=12.1 -c pytorch -c nvidia
# On Mac
# conda install -y pytorch::pytorch torchaudio -c pytorch

conda install -y pysoundfile -c conda-forge

conda install -y ffmpeg-python -c conda-forge

pip install sox

pip install uroman # conda does not have it
cp /opt/conda/envs/mms_fa/bin/uroman /opt/conda/envs/mms_fa/bin/uroman.pl
# on Mac
# cp /Users/gary/miniforge3/envs/mms_fa/bin/uroman /Users/gary/miniforge3/envs/mms_fa/bin/uroman.pl

conda deactivate

# ---------- whisper ----------

# https://pypi.org/project/openai-whisper/

conda create -y -n whisper python=3.11

conda activate whisper

pip install -U openai-whisper

conda deactivate

# ---------- go ----------

# install go
#sudo yum -y install golang
wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
export PATH=/usr/local/go/bin:$PATH
rm go1.23.3.linux-amd64.tar.gz

# Check go
which go
go version
go env

# Set GOPATH
cd $HOME
mkdir go
export GOPATH=$HOME/go

# ---------- application server ----------

mkdir $HOME/go/src
cd $HOME/go/src
git clone https://github.com/garygriswold/fcbh-dataset-io.git
mv fcbh-dataset-io dataset  # Fix this hack

# Update and compile
# cd $HOME/go/src/dataset
# get pull origin main
# go install dataset/controller/api_server

cd $HOME
mkdir data
mkdir data/download
mkdir data/tmp

# ---------- environment variables ----------

FILE=/home/ec2-user/.bashrc
echo "# >>> start DATASET initialize >>>" >> $FILE
echo "# GOPATH" >> $FILE
echo "export GOPATH=$HOME/go" >> $FILE
echo "export GOPROJ=$GOPATH/src" >> $FILE
echo "export PATH=$PATH:/usr/local/go/bin" >> $FILE
echo "# OPEN_AI" >> $FILE
echo "export OPENAI_API_KEY=<redacted>" >> $FILE
echo "# Bible Brain Key" >> $FILE
echo "export FCBH_DBP_KEY=<redacted>" >> $FILE
echo "# Data" >> $FILE
echo "export FCBH_DATASET_DB=$HOME/data" >> $FILE
echo "export FCBH_DATASET_FILES=$HOME/data/download" >> $FILE
echo "export FCBH_DATASET_TMP=$HOME/data/tmp" >> $FILE
echo "export FCBH_DATASET_QUEUE=dataset-queue" >> $FILE
echo "export FCBH_DATASET_IO_BUCKET=dataset-io" >> $FILE
echo "# Conda ENVS" >> $FILE
echo "export FCBH_AENEAS_PYTHON=/opt/conda/envs/aeneas/bin/python" >> $FILE
echo "export FCBH_FASTTEXT_EXE=/opt/conda/envs/fasttext/fastText/fasttext" >> $FILE
echo "export FCBH_LIBROSA_PYTHON=/opt/conda/envs/librosa/bin/python" >> $FILE
echo "export FCBH_MMS_ASR_PYTHON=/opt/conda/envs/mms_asr/bin/python" >> $FILE
echo "export FCBH_MMS_FA_PYTHON=/opt/conda/envs/mms_fa/bin/python" >> $FILE
echo "export FCBH_UROMAN_EXE=/opt/conda/envs/mms_fa/bin/uroman" >> $FILE
echo "export FCBH_WHISPER_EXE=/opt/conda/envs/whisper/bin/whisper" >> $FILE
echo "# Logging" >> $FILE
echo "export FCBH_DATASET_LOG_FILE=$HOME/dataset.log" >> $FILE
echo "export FCBH_DATASET_LOG_LEVEL=DEBUG" >> $FILE
echo "export PYTORCH_CUDA_ALLOC_CONF=expandable_segments:True" >> $FILE
echo "# <<< end DATASET initialize <<<" >> $FILE

source /home/ec2-user/.bashrc

# ---------- aws credentials ----------

cd /home/ec2-user
mkdir .aws
FILE=.aws/credentials
echo "[default]" > $FILE
echo "aws_access_key_id = <redacted>" >> $FILE
echo "aws_secret_access_key = <redacted>" >> $FILE

# ---------- system ctl file ----------

##
## Update these env values after CLI is tested.
##
FILE=$HOME/tmp.service
echo "[Unit]" > $FILE
echo 'Description="FCBH Dataset API Server"' >> $FILE
echo "After=network.target" >> $FILE
echo "" >> $FILE
echo "[Service]" >> $FILE
echo "Type=simple" >> $FILE
echo "User=dataset" >> $FILE
echo "# GOPATH" >> $FILE
echo 'Environment="GOPATH=$HOME/go"' >> $FILE
echo 'Environment="GOPROJ=$GOPATH/src"' >> $FILE
echo 'Environment="FCBH_DBP_KEY=<redacted>"' >> $FILE
echo 'Environment="FCBH_DATASET_DB=$HOME/data"' >> $FILE
echo 'Environment="FCBH_DATASET_FILES=$HOME/data/download"' >> $FILE
echo 'Environment="FCBH_DATASET_TMP=$HOME/data/tmp"' >> $FILE
echo 'Environment="PYTHON_EXE=/opt/conda/envs/aeneas/bin/python"' >> $FILE
echo 'Environment="FCBH_AENEAS_PYTHON=/opt/conda/envs/aeneas/bin/python"' >> $FILE
echo 'Environment="FCBH_MMS_PYTHON=/opt/conda/envs/mms_fa/bin/python"' >> $FILE
echo 'Environment="WHISPER_EXE=/opt/conda/envs/aeneas/bin/whisper"' >> $FILE
echo 'Environment="FCBH_DATASET_LOG_FILE=/home/dataset/dataset.log"' >> $FILE
echo 'Environment="FCBH_DATASET_LOG_LEVEL=DEBUG"' >> $FILE
echo "ExecStart=/home/dataset/go/bin/api_server" >> $FILE
echo "ExecReload=/bin/kill -HUP $MAINPID" >> $FILE
echo "Restart=always" >> $FILE
echo "RestartSec=5" >> $FILE
echo "" >> $FILE
echo "[Install]" >> $FILE
echo "WantedBy=multi-user.target" >> $FILE
echo "" >> $FILE
echo "# sudo systemctl daemon-reload" >> $FILE
echo "" >> $FILE
echo "# sudo systemctl enable api_server" >> $FILE
echo "# sudo systemctl start api_server" >> $FILE
echo "# sudo systemctl stop api_server" >> $FILE
echo "# sudo systemctl restart api_server" >> $FILE
echo "# sudo systemctl reload api_server" >> $FILE
echo "# sudo systemctl status api_server" >> $FILE
echo "# sudo systemctl disable api_server" >> $FILE

sudo cp $FILE /etc/systemd/system/fcbh_dataset_io.service

# End of ec2-user commands
EOF

# Make sure all files are owned by ec2-user
# Set proper ownership of all installed files
chown -R ec2-user:ec2-user /opt/conda
chown -R ec2-user:ec2-user /home/ec2-user/.conda

# Set proper permissions
chmod -R u+rwX,go-w /opt/conda

sudo -u ec2-user bash << 'EOF2'

source /home/ec2-user/.bashrc
cd $GOPROJ/dataset/
go test ./tests -run TestTSBB

go install dataset/controller/queue_server

#/home/ec2-user/go/bin/queue_server &


EOF2