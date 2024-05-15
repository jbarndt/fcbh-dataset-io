# Installation Instructions for FCBH IO Dataset Server

## Linux Distro -> Centos

## Utilities to install

install sqlite3

Required for Aeneas:
ffmpeg 
ffprobe (usually provided by the ffmpeg package)
espeak 

Make sure python3 is installed

$ sudo pip install numpy
$ sudo pip install aeneas

$ pip install librosa

> $ pip install -U openai-whisper

`$ pip install fasttext`

$ go
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version

$ git
$ sudo dnf install git-all

Install Application
> $ git clone https://github.com/garygriswold/fcbh-dataset-io.git