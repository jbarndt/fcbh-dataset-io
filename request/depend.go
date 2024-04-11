package request

import (
	"dataset"
	"strings"
)

func (r *RequestDecoder) Prereq(req *Request) {
	mfccPrereq(req)
}

func mfccPrereq(req *Request) {
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
