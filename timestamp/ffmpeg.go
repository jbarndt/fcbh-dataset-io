package timestamp

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	log "dataset/logger"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ChopByTimestamp uses timestamps to chop timestamps into files, and puts the filenames in timestamp record.
func ChopByTimestamp(ctx context.Context, tempDir string, file input.InputFile, timestamps []db.Audio) ([]db.Audio, dataset.Status) {
	var results []db.Audio
	var status dataset.Status
	var command []string
	command = append(command, `-i`, file.FilePath())
	command = append(command, `-codec:a`, `copy`)
	command = append(command, `-y`)
	for _, ts := range timestamps {
		if ts.BeginTS == 0.0 && ts.EndTS == 0.0 {
			continue
		}
		beginTS := strconv.FormatFloat(ts.BeginTS, 'f', 2, 64)
		command = append(command, `-ss`, beginTS)
		if ts.EndTS != 0.0 {
			endTS := strconv.FormatFloat(ts.EndTS, 'f', 2, 64)
			command = append(command, `-to`, endTS)
		}
		verseFilename := fmt.Sprintf("verse_%s_%d_%s_%s.wav",
			file.BookId, file.Chapter, ts.VerseStr, beginTS)
		ts.AudioVerse = filepath.Join(tempDir, verseFilename)
		command = append(command, `-c`, `copy`, ts.AudioVerse)
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

// ConvertMp3toWav
func ConvertMp3ToWav(ctx context.Context, tempDir string, file input.InputFile) (string, dataset.Status) {
	// ffmpeg -I filename.mp3 -acodec pcm_s16le -ar 16000 output.wav
	var outputPath string
	var status dataset.Status
	if filepath.Ext(file.Filename) == ".wav" {
		outputPath = file.FilePath()
	} else {
		filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
		outputPath = filepath.Join(tempDir, filename+".wav")
		ffMpegPath := `ffmpeg`
		cmd := exec.Command(ffMpegPath,
			file.FilePath(),
			`-acodec`, `pcm_s16le`,
			`-ar`, `16000`,
			outputPath)
		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf
		err := cmd.Run()
		if err != nil {
			status = log.Error(ctx, 500, err, stderrBuf.String())
		}
	}
	return outputPath, status
}
