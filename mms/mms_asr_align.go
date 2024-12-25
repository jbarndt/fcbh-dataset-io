package mms

import (
	"bufio"
	"dataset"
	"dataset/generic"
	log "dataset/logger"
	"dataset/timestamp"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ProcessAlignSilence will perform Auto Speech Recognition on the silent parts audio gaps
func (a *MMSASR) ProcessAlignSilence(directory string, chars []generic.AlignChar) ([]generic.AlignChar, dataset.Status) {
	var result []generic.AlignChar
	lang, status := checkLanguage(a.ctx, a.lang, a.sttLang, "mms_asr")
	if status.IsErr {
		return result, status
	}
	pythonScript := filepath.Join(os.Getenv("GOPROJ"), "dataset/mms/mms_asr.py")
	writer, reader, status := callStdIOScript(a.ctx, os.Getenv(`FCBH_MMS_ASR_PYTHON`), pythonScript, lang)
	if status.IsErr {
		return result, status
	}
	tempDir, err := os.MkdirTemp(os.Getenv(`FCBH_DATASET_TMP`), "mms_asr_align_")
	if err != nil {
		return result, log.Error(a.ctx, 500, err, `Error creating temp dir`)
	}
	defer os.RemoveAll(tempDir)
	words := a.groupByWord(chars)
	var wavFile string
	var priorAudioFile string
	for _, word := range words {
		var silentChar generic.AlignChar
		for _, char := range word {
			if char.SilenceLong > 0 {
				silentChar = char
			}
			result = append(result, char)
		}
		if silentChar.LineRef != "" {
			firstCh := word[0]
			lastCh := word[len(word)-1]
			beginTS := firstCh.BeginTS                     // - 0.025???
			endTS := silentChar.EndTS + silentChar.Silence // - 0.025
			if endTS < lastCh.EndTS {
				endTS = lastCh.EndTS
			}
			s := silentChar
			if beginTS >= endTS {
				fmt.Println(s.LineRef, s.WordId, s.SilencePos, s.SilenceLong,
					math.Round(s.Silence*100.0)/100.0, "Not enough time to align silence")
				continue
			}
			if priorAudioFile != s.AudioFile {
				priorAudioFile = s.AudioFile
				filePath := filepath.Join(directory, s.AudioFile)
				wavFile, status = timestamp.ConvertMp3ToWav(a.ctx, tempDir, filePath)
				if status.IsErr {
					return result, status
				}
				resp, stat := a.asrFile(writer, reader, wavFile)
				if stat.IsErr {
					return result, status
				}
				fmt.Println(resp, "\n")
			}
			audioFragment, status3 := timestamp.ChopOneSegment(a.ctx, tempDir, wavFile, beginTS, endTS)
			if status3.IsErr {
				return result, status3
			}
			response, status4 := a.asrFile(writer, reader, audioFragment)
			if status4.IsErr {
				return result, status4
			}
			//_, err = writer.WriteString(audioFragment + "\n")
			//if err != nil {
			//	return chars, log.Error(a.ctx, 500, err, "Error writing to mms_asr.py")
			//}
			//err = writer.Flush()
			//if err != nil {
			//	return chars, log.Error(a.ctx, 500, err, "Error flush to mms_asr.py")
			//}
			//response, err2 := reader.ReadString('\n')
			//if err2 != nil {
			//	return chars, log.Error(a.ctx, 500, err2, `Error reading mms_asr.py response`)
			//}
			//response = strings.TrimRight(response, "\n")
			var str []rune
			for _, c := range word {
				str = append(str, c.CharNorm)
			}
			fmt.Println(s.LineRef, s.WordId, s.SilencePos, s.SilenceLong,
				math.Round(s.Silence*100.0)/100.0, string(str), response)
			for _, resp := range response {
				var newChar generic.AlignChar
				newChar.AudioFile = s.AudioFile
				newChar.LineId = s.LineId
				newChar.LineRef = s.LineRef
				//newChar.WordId = ch.WordId // might not be correct
				newChar.CharNorm = resp
				newChar.BeginTS = beginTS
				newChar.EndTS = endTS // It a number of chars are found they have the same TS
				newChar.FAScore = 1.0
				newChar.IsASR = true
				result = append(result, newChar)
			}
		}
	}
	log.Debug(a.ctx, "Finished ASR Align")
	bytes, err := json.Marshal(chars)
	if err != nil {
		log.Warn(a.ctx, 500, err, "Error creating json file of ASR Align result")
	}
	// This is supposed to be stored with the database
	_ = os.WriteFile("mms_asr_align.json", bytes, 0644)
	return chars, status
}

func (a *MMSASR) groupByWord(chars []generic.AlignChar) [][]generic.AlignChar {
	var result [][]generic.AlignChar
	if len(chars) == 0 {
		return result
	}
	currWordId := chars[0].WordId
	start := 0
	for i, ch := range chars {
		if currWordId != ch.WordId {
			currWordId = ch.WordId
			oneLine := make([]generic.AlignChar, i-start)
			copy(oneLine, chars[start:i])
			start = i
			result = append(result, oneLine)
		}
	}
	return result
}

func (a *MMSASR) computeMidPoints(i int, words [][]generic.AlignChar) (float64, float64) {
	var beginTS float64
	var endTS float64
	if i == 0 {
		currWord := words[i]
		nextWord := words[i+1]
		beginTS = 0.0
		endTS = (currWord[len(currWord)-1].EndTS + nextWord[0].BeginTS) / 2.0
	} else if i < len(words) {
		priorWord := words[i-1]
		currWord := words[i]
		nextWord := words[i+1]
		beginTS = (priorWord[len(priorWord)-1].EndTS + currWord[0].BeginTS) / 2.0
		endTS = (currWord[len(currWord)-1].EndTS + nextWord[0].BeginTS) / 2.0
	} else {
		priorWord := words[i-1]
		currWord := words[i]
		beginTS = (priorWord[len(priorWord)-1].EndTS + currWord[0].BeginTS) / 2.0
		endTS = currWord[len(currWord)-1].EndTS + 0.5
	}
	return beginTS, endTS
}

func (a *MMSASR) asrFile(writer *bufio.Writer, reader *bufio.Reader, audioFile string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	_, err := writer.WriteString(audioFile + "\n")
	if err != nil {
		return result, log.Error(a.ctx, 500, err, "Error writing to mms_asr.py")
	}
	err = writer.Flush()
	if err != nil {
		return result, log.Error(a.ctx, 500, err, "Error flush to mms_asr.py")
	}
	response, err2 := reader.ReadString('\n')
	if err2 != nil {
		return result, log.Error(a.ctx, 500, err2, `Error reading mms_asr.py response`)
	}
	response = strings.TrimRight(response, "\n")
	return response, status
}
