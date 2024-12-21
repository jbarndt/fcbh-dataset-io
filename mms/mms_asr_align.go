package mms

import (
	"dataset"
	"dataset/generic"
	log "dataset/logger"
	"dataset/timestamp"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProcessAlignSilence will perform Auto Speech Recognition on the silent parts audio gaps
func (a *MMSASR) ProcessAlignSilence(directory string, lines []generic.AlignLine) ([]generic.AlignLine, dataset.Status) {
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return lines, status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_asr.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_ASR_PYTHON`), pythonScript, lang)
	if status.IsErr {
		return lines, status
	}
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_asr_align_")
	if err != nil {
		return lines, log.Error(a.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	var wavFile string
	var priorAudioFile string
	for i, line := range lines {
		var newLine []generic.AlignChar
		for _, ch := range line.Chars {
			newLine = append(newLine, ch)
			if ch.SilenceLong != 0 {
				if priorAudioFile != ch.AudioFile {
					priorAudioFile = ch.AudioFile
					filePath := filepath.Join(directory, ch.AudioFile)
					wavFile, status = timestamp.ConvertMp3ToWav(a.ctx, tempDir, filePath)
					if status.IsErr {
						return lines, status
					}
				}
				endSilenceTS := ch.EndTS + ch.Silence
				audioFragment, status3 := timestamp.ChopOneSegment(a.ctx, tempDir, wavFile, ch.EndTS, endSilenceTS)
				if status3.IsErr {
					return lines, status3
				}
				_, err = writer.WriteString(audioFragment + "\n")
				if err != nil {
					return lines, log.Error(a.ctx, 500, err, "Error writing to mms_asr.py")
				}
				err = writer.Flush()
				if err != nil {
					return lines, log.Error(a.ctx, 500, err, "Error flush to mms_asr.py")
				}
				response, err2 := reader.ReadString('\n')
				if err2 != nil {
					return lines, log.Error(a.ctx, 500, err2, `Error reading mms_asr.py response`)
				}
				response = strings.TrimRight(response, "\n")
				fmt.Println(ch.LineRef, ch.WordId, response)
				for _, resp := range response {
					var newChar generic.AlignChar
					newChar.AudioFile = ch.AudioFile
					newChar.LineId = ch.LineId
					newChar.LineRef = ch.LineRef
					//newChar.WordId = ch.WordId // might not be correct
					newChar.CharNorm = resp
					newChar.BeginTS = ch.EndTS
					newChar.EndTS = endSilenceTS // It a number of chars are found they have the same TS
					newChar.FAScore = 1.0
					newChar.IsASR = true
					newLine = append(newLine, newChar)
				}
			}
		}
		lines[i].Chars = newLine
	}
	log.Debug(a.ctx, "Finished ASR Align")
	bytes, err := json.Marshal(lines)
	if err != nil {
		log.Warn(a.ctx, 500, err, "Error creating json file of ASR Align result")
	}
	// This is supposed to be stored with the database
	_ = os.WriteFile("mms_asr_align.json", bytes, 0644)
	return lines, status
}
