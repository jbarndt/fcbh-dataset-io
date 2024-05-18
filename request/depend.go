package request

import (
	"dataset"
	"path/filepath"
	"strings"
)

func (r *RequestDecoder) Prereq(req *Request) {
	r.mfccPrereq(req)
	r.setOutputType(req)
}

func (r *RequestDecoder) mfccPrereq(req *Request) {
	if req.AudioEncoding.MFCC {
		if req.Timestamps.NoTimestamps {
			if !req.AudioData.NoAudio {
				if req.TextData.NoText {
					req.Timestamps.Aeneas = true
				} else {
					req.Timestamps.BibleBrain = true
				}
			}
		}
	}
}

func (r *RequestDecoder) setOutputType(req *Request) dataset.Status {
	var status dataset.Status
	var msgs []string
	fType := strings.ToLower(filepath.Ext(req.OutputFile))
	switch fType {
	case ".json":
		req.OutputFormat.JSON = true
	case ".csv":
		req.OutputFormat.CSV = true
	case ".sqlite":
		req.OutputFormat.Sqlite = true
	case ".html":
		req.OutputFormat.HTML = true
	default:
		msg := `output_file must be .json, .csv, .sqlite, or .html for compare tasks`
		msgs = append(msgs, msg)
	}
	return status
}

func (r *RequestDecoder) Depend(req Request) dataset.Status {
	var status dataset.Status
	var msgs []string
	if !req.Timestamps.NoTimestamps {
		if req.AudioData.NoAudio {
			msg := `Timestamps are requested, but there is no audio`
			msgs = append(msgs, msg)
		}
	}
	if req.Timestamps.Aeneas {
		if req.TextData.NoText {
			msg := `Aeneas requested, but there is no text data`
			msgs = append(msgs, msg)
		}
	}
	if !req.TextEncoding.NoEncoding {
		if req.TextData.NoText {
			msg := `Text encoding requested, but there is no text data`
			msgs = append(msgs, msg)
		}
	}
	if len(msgs) > 0 {
		status.IsErr = true
		status.Status = 400
		status.Message = strings.Join(msgs, "\n")
	}
	return status
}
