# Installation Instructions for FCBH IO Dataset Server

> Create server Ubuntu 20.04 x64

167.99.58.202 ubuntu 20

> adduser dataset

> usermod -aG sudo dataset

exit & login dataset

> sudo apt update

> sudo apt -y upgrade

> python3 -V

3.8.10

> sudo apt install ffmpeg

> sudo apt install espeak

> sudo apt install python3-pip

> pip install numpy

> source .profile  # needed to add $HOME/.local/bin

> pip install aeneas -> failed

> wget https://raw.githubusercontent.com/readbeyond/aeneas/master/install_dependencies.sh

> bash install_dependencies.sh

> pip install aeneas

> pip install librosa

> pip install -U openai-whisper ## This failed because of memory

Because whisper failed to install, these things were tried
based on information on the whisper page.
- sudo apt install rust-all
- pip install tiktoken
- export PATH="$HOME/.cargo/bin:$PATH"
- pip install setuptools-rust
+ It is not clear if they were needed

> pip install torch --no-cache-dir

> pip install -U openai-whisper

> pip install fasttext

> sudo apt install sqlite3

Install go
> wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz

> sudo rm -rf /usr/local/go

> sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz

> vi .profile
insert PATH=$PATH:/usr/local/go/bin

Install Application Server
> cd $HOME/go/src

> git clone https://github.com/garygriswold/fcbh-dataset-io.git

> mv fcbh-dataset-io dataset  # Fix this hack

> cd dataset

> go run dataset/controller/api_server

add environment variables
 
> export FCBH_DBP_KEY=b4xxx
 
> export FCBH_DATASET_DB=$HOME/data
 
> export FCBH_DATASET_FILES=$HOME/data/download
 
> export FCBH_DATASET_TMP=$HOME/data/tmp 

> export PYTHON_EXE=/usr/bin/python3

> export WHISPER_EXE=$HOME/.local/bin/whisper

> cd $HOME/go/src/dataset

> go install dataset/controller/api_server

> cd $HOME

> ./go/bin/api_server &

To Update go code

> cd $HOME/go/src/dataset

> get pull origin main

> go install dataset/controller/api_server

> sudo vi /etc/systemd/system/api_server.service

insert: dataset/controller/api_server/api_server.service

> sudo systemctl start api_server

> sudo systemctl status api_server

https://github.com/readbeyond/aeneas/blob/master/wiki/INSTALL.md

https://github.com/openai/whisper