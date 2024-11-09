package encode

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"dataset/request"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type MFCC struct {
	ctx     context.Context
	conn    db.DBAdapter
	bibleId string
	detail  request.Detail
	numMFCC int
}

func NewMFCC(ctx context.Context, conn db.DBAdapter, bibleId string,
	detail request.Detail, numMFCC int) MFCC {
	var m MFCC
	m.ctx = ctx
	m.conn = conn
	m.bibleId = bibleId
	m.detail = detail
	m.numMFCC = numMFCC
	return m
}

func (m *MFCC) ProcessFiles(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	for _, aFile := range audioFiles {
		var mfccResp MFCCResp
		mfccResp, status = m.executeLibrosa(aFile.FilePath())
		if status.IsErr {
			return status
		}
		if m.detail.Lines {
			status = m.processScripts(mfccResp, aFile.BookId, aFile.Chapter)
			if status.IsErr {
				return status
			}
		}
		if m.detail.Words {
			status = m.processWords(mfccResp, aFile.BookId, aFile.Chapter)
			if status.IsErr {
				return status
			}
		}
	}
	return status
}

type MFCCResp struct {
	AudioFile  string      `json:"input_file"`
	SampleRate float64     `json:"sample_rate"`
	HopLength  float64     `json:"hop_length"`
	FrameRate  float64     `json:"frame_rate"`
	Shape      []int       `json:"mfcc_shape"`
	Type       string      `json:"mfcc_type"`
	MFCC       [][]float32 `json:"mfccs"`
}

func (m *MFCC) executeLibrosa(audioFile string) (MFCCResp, dataset.Status) {
	var result MFCCResp
	var status dataset.Status
	pythonPath := os.Getenv(`PYTHON_EXE`)
	mfccLibrosaPath := filepath.Join(os.Getenv(`GOPROJ`), `dataset`, `encode`, `mfcc_librosa.py`)
	cmd := exec.Command(pythonPath, mfccLibrosaPath, audioFile, strconv.Itoa(m.numMFCC))
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(m.ctx, 500, err, `Error executing mfcc_librosa.py`)
		// Don't return here, need to see stderr
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

func (m *MFCC) processScripts(mfcc MFCCResp, bookId string, chapterNum int) dataset.Status {
	var status dataset.Status
	timestamps, status := m.conn.SelectScriptTimestamps(bookId, chapterNum)
	mfccs := m.segmentMFCC(timestamps, mfcc)
	status = m.conn.InsertScriptMFCCS(mfccs)
	return status
}

func (m *MFCC) processWords(mfcc MFCCResp, bookId string, chapterNum int) dataset.Status {
	var status dataset.Status
	timestamps, status := m.conn.SelectWordTimestamps(bookId, chapterNum)
	mfccs := m.segmentMFCC(timestamps, mfcc)
	status = m.conn.InsertWordMFCCS(mfccs)
	return status
}

func (m *MFCC) segmentMFCC(timestamps []db.Timestamp, mfcc MFCCResp) []db.MFCC {
	var result []db.MFCC
	for _, ts := range timestamps {
		startIndex := int(ts.BeginTS*mfcc.FrameRate + 0.5)
		endIndex := int(ts.EndTS*mfcc.FrameRate + 0.5)
		var segment [][]float32
		if endIndex != 0 {
			segment = mfcc.MFCC[startIndex:endIndex][:]
		} else {
			segment = mfcc.MFCC[startIndex:][:] // The last timestamp from DB will be zero
		}
		var mf db.MFCC
		mf.Id = ts.Id
		mf.Rows = len(segment)
		if mf.Rows == 0 {
			mf.Cols = 0
		} else {
			mf.Cols = len(segment[0])
		}
		mf.MFCC = segment
		result = append(result, mf)
	}
	return result
}
