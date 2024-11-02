package request

import (
	"path/filepath"
	"strings"
)

func (r *RequestDecoder) Prereq(req *Request) {
	//	r.mfccPrereq(req)
	r.setOutputType(req)
}

//func (r *RequestDecoder) mfccPrereq(req *Request) {
//	if req.AudioEncoding.MFCC {
//		if req.Timestamps.NoTimestamps {
//			if !req.AudioData.NoAudio {
//				if req.TextData.NoText {
//					req.Timestamps.Aeneas = true
//				} else {
//					req.Timestamps.BibleBrain = true
//				}
//			}
//		}
//	}
//}

func (r *RequestDecoder) setOutputType(req *Request) {
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
		r.errors = append(r.errors, msg)
	}
}

func (r *RequestDecoder) Depend(req Request) {
	if !req.Timestamps.NoTimestamps {
		if req.AudioData.NoAudio {
			r.errors = append(r.errors, `Timestamps are requested, but there is no audio`)
		}
	}
	if req.Timestamps.Aeneas || req.Timestamps.MMSFAVerse || req.Timestamps.MMSFAWord {
		if req.TextData.NoText {
			r.errors = append(r.errors, `Timestamp estimation requested, but there is no text data`)
		}
	}
	if !req.TextEncoding.NoEncoding {
		if req.TextData.NoText {
			r.errors = append(r.errors, `Text encoding requested, but there is no text data`)
		}
	}
	if !req.SpeechToText.NoSpeechToText {
		if req.AudioData.NoAudio {
			r.errors = append(r.errors, `Speech to Text is requested, but there is no audio`)
		}
		if req.Timestamps.NoTimestamps {
			r.errors = append(r.errors, `Speech to Text is requested, but there are no timestamps`)
		}
	}
	if req.AudioEncoding.MFCC {
		if req.Timestamps.NoTimestamps {
			r.errors = append(r.errors, `MFCC's are requested', but there are no timestamps`)
		}
	}
}
