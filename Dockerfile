FROM public.ecr.aws/docker/library/golang:1.23.4-alpine as builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o taskService controller/ecs_entry.go

# Deployment stage
FROM public.ecr.aws/amazonlinux/amazonlinux:latest

WORKDIR /root/


# Install necessary tools
RUN yum update -y && \
    # yum install -y wget tar xz && \
    yum clean all




FROM ubuntu:20.04

# Install base utilities
RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get install -y wget \
    && apt-get install -y git \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y cmake \
    && apt-get install -y build-essential


# Install conda
ENV CONDA_DIR /opt/conda
RUN wget --quiet https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O ~/miniconda.sh && \
    /bin/bash ~/miniconda.sh -b -p /opt/conda

# Put conda in path 
ENV PATH=$CONDA_DIR/bin:$PATH

WORKDIR /app

# ---------- base ----------
COPY environment_base.yml .
RUN conda env update -f environment_base.yml


# ---------- aeneas ----------
COPY environment_aeneas.yml .
RUN conda env create -f environment_aeneas.yml
RUN conda run -n aeneas apt-get install -y espeak libespeak-dev build-essential
RUN conda run -n aeneas pip install numpy
RUN conda run -n aeneas pip install aeneas


# ---------- easy_mms ----------
COPY environment_easy_mms.yml .
RUN conda env create -f environment_easy_mms.yml
#RUN conda run -n easy_mms pip install easymms    TODO - figure out error


# ---------- fasttext ----------
RUN conda create -y -n fasttext --no-default-packages
RUN git clone https://github.com/facebookresearch/fastText.git
RUN mkdir build
WORKDIR /app/fastText/build
RUN cmake ..
RUN make && make install 
WORKDIR /app


# ---------- librosa ----------
COPY environment_librosa.yml .
RUN conda env create -f environment_librosa.yml
RUN conda run -n librosa pip install librosa


# ---------- mms_asr ----------
COPY environment_mms_asr.yml .
RUN conda env create -f environment_mms_asr.yml
RUN conda run -n mms_asr pip install accelerate datasets soundfile librosa 
RUN conda run -n mms_asr pip install --upgrade transforms
RUN conda run -n mms_asr pip install uroman
RUN cp /opt/conda/envs/mms_asr/bin/uroman /opt/conda/envs/mms_asr/bin/uroman.pl


# ---------- mms_fa ----------
COPY environment_mms_fa.yml .
RUN conda env create -f environment_mms_fa.yml
RUN conda run -n mms_fa pip install sox uroman
RUN cp /opt/conda/envs/mms_fa/bin/uroman /opt/conda/envs/mms_fa/bin/uroman.pl


# ---------- whisper ----------
COPY environment_whisper.yml .
RUN conda env create -f environment_whisper.yml
RUN conda run -n whisper pip install -U openai-whisper


# ---------- go ----------
RUN wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH
RUN rm go1.23.3.linux-amd64.tar.gz
RUN mkdir /app/go
ENV GOPATH=/app/go/


# ---------- application server ----------
RUN mkdir /app/go/src
WORKDIR /app/go/src
RUN git clone https://github.com/garygriswold/fcbh-dataset-io.git
RUN mv fcbh-dataset-io dataset

WORKDIR /app
RUN mkdir /app/data && mkdir /app/data/download && mkdir /app/data/tmp


# ---------- environment variables ----------
ENV HOME=/app
ENV GOPROJ=$GOPATH/src
ENV PATH=$PATH:/usr/local/go/bin
ENV FCBH_DATASET_DB=$HOME/data
ENV FCBH_DATASET_FILES=$HOME/data/download
ENV FCBH_DATASET_TMP=$HOME/data/tmp
ENV FCBH_AENEAS_PYTHON=/opt/conda/envs/aeneas/bin/python
ENV FCBH_FASTTEXT_EXE=/opt/conda/envs/fasttext/fastText/fasttext
ENV FCBH_LIBROSA_PYTHON=/opt/conda/envs/librosa/bin/python
ENV FCBH_MMS_ASR_PYTHON=/opt/conda/envs/mms_asr/bin/python
ENV FCBH_MMS_FA_PYTHON=/opt/conda/envs/mms_fa/bin/python
ENV FCBH_UROMAN_EXE=/opt/conda/envs/mms_fa/bin/uroman
ENV FCBH_WHISPER_EXE=/opt/conda/envs/whisper/bin/whisper
ENV FCBH_DATASET_LOG_FILE=$HOME/dataset.log
ENV FCBH_DATASET_LOG_LEVEL=DEBUG
ENV PYTORCH_CUDA_ALLOC_CONF=expandable_segments:True




# Copy the built binary
COPY --from=builder /app/taskService .

# Ensure the binary is executable
RUN chmod +x ./taskService

# Use environment variable to determine the service to run
ENTRYPOINT ["./taskService"]
