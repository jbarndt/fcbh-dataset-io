package timestamp

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ChopByTimestamp uses timestamps to chop timestamps into files, and puts the filenames in timestamp record.
func ChopByTimestamp(ctx context.Context, tempDir string, inputFile string, timestamps []db.Audio) ([]db.Audio, dataset.Status) {
	var results []db.Audio
	var status dataset.Status
	var fileExt = filepath.Ext(inputFile)
	var command []string
	command = append(command, `-i`, inputFile)
	command = append(command, `-codec:a`, `copy`)
	command = append(command, `-y`)
	for _, ts := range timestamps {
		if ts.BeginTS == 0.0 && ts.EndTS == 0.0 {
			continue
		}
		beginTS := strconv.FormatFloat(ts.BeginTS, 'f', 3, 64)
		command = append(command, `-ss`, beginTS)
		if ts.EndTS != 0.0 {
			endTS := strconv.FormatFloat(ts.EndTS, 'f', 3, 64)
			command = append(command, `-to`, endTS)
		}
		verseFilename := fmt.Sprintf("verse_%s_%d_%s_%s%s",
			ts.BookId, ts.ChapterNum, ts.VerseStr, beginTS, fileExt)
		ts.AudioVerseWav = filepath.Join(tempDir, verseFilename)
		command = append(command, `-c`, `copy`, ts.AudioVerseWav)
		results = append(results, ts)
	}
	ffMpegPath := `ffmpeg`
	cmd := exec.Command(ffMpegPath, command...)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(ctx, 500, err, stderrBuf.String())
	}
	return results, status
}

// ChopOneSegment uses timestamps extract one segment from an audio file
func ChopOneSegment(ctx context.Context, tempDir string, inputFile string, beginTS float64, endTS float64) (string, dataset.Status) {
	var outputFile string
	var status dataset.Status
	outputFile = filepath.Join(tempDir, fmt.Sprintf("%d.wav", time.Now().UnixNano()))
	err := ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs{
		"codec:a": "copy",
		"c":       "copy",
		"y":       "",
		"ss":      beginTS,
		"to":      endTS,
	}).OverWriteOutput().Run()
	if err != nil {
		status = log.Error(ctx, 500, err, "Error in ChopOneSegment")
	}
	return outputFile, status
}

func ConvertMp3ToWav(ctx context.Context, tempDir string, inputFile string) (string, dataset.Status) {
	var outputPath string
	var status dataset.Status
	filename := filepath.Base(inputFile)
	outputFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputPath = filepath.Join(tempDir, outputFilename+".wav")
	err := ffmpeg.Input(inputFile).Output(outputPath, ffmpeg.KwArgs{
		"acodec": "pcm_s16le",
		"ar":     "16000",
		"ac":     "1",
	}).OverWriteOutput().Run()
	if err != nil {
		status = log.Error(ctx, 500, err, "Error ")
	}
	return outputPath, status
}

// ConvertMp3toWav
func OldConvertMp3ToWav(ctx context.Context, tempDir string, filePath string) (string, dataset.Status) {
	// ffmpeg -i filename.mp3 -acodec pcm_s16le -ar 16000 output.wav
	var outputPath string
	var status dataset.Status
	filename := filepath.Base(filePath)
	outputFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputPath = filepath.Join(tempDir, outputFilename+".wav")
	ffMpegPath := `ffmpeg`
	cmd := exec.Command(ffMpegPath,
		`-i`, filePath,
		`-acodec`, `pcm_s16le`,
		`-ar`, `16000`,
		`-ac`, `1`,
		outputPath)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(ctx, 500, err, stderrBuf.String())
	}
	return outputPath, status
}
