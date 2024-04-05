package encode

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type MFCC struct {
	ctx       context.Context
	conn      db.DBAdapter
	bibleId   string
	audioFSId string
}

func NewMFCC(ctx context.Context, conn db.DBAdapter, bibleId string, audioFSId string) MFCC {
	var m MFCC
	m.ctx = ctx
	m.conn = conn
	m.bibleId = bibleId
	m.audioFSId = audioFSId
	return m
}

func (m *MFCC) Process(detail dataset.TextDetailType) dataset.Status {
	var status dataset.Status
	audioFiles, status := ReadDirectory(m.ctx, m.bibleId, m.audioFSId)
	if status.IsErr {
		return status
	}
	if detail == dataset.LINES || detail == dataset.BOTH {
		status = m.processScripts(audioFiles)
	} else if detail == dataset.WORDS || detail == dataset.BOTH {
		status = m.processWords(audioFiles)
	}
	return status
}

func (m *MFCC) processScripts(audioFiles []string) dataset.Status {
	var status dataset.Status
	for _, audioFile := range audioFiles {
		var mfccResp MFCCResp
		mfccResp, status = m.executeLibrosa(audioFile)
		if status.IsErr {
			return status
		}
		fmt.Println(mfccResp.FrameRate, mfccResp.Type, mfccResp.SampleRate, mfccResp.Shape, mfccResp.MFCC[0][0])
		// execute scripts timestamp query
		// possibly move the data to a generic structure, id, beginTS, endTS
		// break up mfcc data into timestamp segments by ts
		// update script MFCC segments, create a new type

		os.Exit(0)
	}
	return status
}

func (m *MFCC) processWords(audioFiles []string) dataset.Status {
	var status dataset.Status
	return status
}

type MFCCResp struct {
	AudioFile  string      `json:"input_file"`
	SampleRate float64     `json:"sample_rate"`
	HopLength  float64     `json:"hop_length"`
	FrameRate  float64     `json:"frame_rate"`
	Shape      []int       `json:"mfcc_shape"`
	Type       string      `json:"mfcc_type"`
	MFCC       [][]float64 `json:"mfccs"`
}

func (m *MFCC) executeLibrosa(audioFile string) (MFCCResp, dataset.Status) {
	var result MFCCResp
	var status dataset.Status
	pythonPath := "python3"
	cmd := exec.Command(pythonPath, `mfcc_librosa.py`, audioFile)
	fmt.Println(cmd.String())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(m.ctx, 500, err, `Error executing mfcc_librosa.py`)
		return result, status
	}
	if stderrBuf.Len() > 0 {
		status = log.ErrorNoErr(m.ctx, 500, `mfcc_librosa.py stderr:`, stderrBuf.String())
		return result, status
	}
	if stdoutBuf.Len() == 0 {
		status = log.ErrorNoErr(m.ctx, 500, `mfcc_librosa.py has no output.`)
		return result, status
	}
	err = json.Unmarshal(stdoutBuf.Bytes(), &result)
	if err != nil {
		status = log.Error(m.ctx, 500, err, `Error parsing json from librosa`)
	}
	return result, status
}
