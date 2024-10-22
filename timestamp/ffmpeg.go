package timestamp

import (
	"bytes"
	"context"
	"dataset"
	"dataset/input"
	log "dataset/logger"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ChopByTimestamp uses timestamps to chop timestamps into files, and puts the filenames in timestamp record.
func ChopByTimestamp(ctx context.Context, tempDir string, file input.InputFile, timestamps []Timestamp) ([]Timestamp, dataset.Status) {
	var results []Timestamp
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
		ts.AudioVerse = fmt.Sprintf("verse_%s_%d_%s_%s.mp3",
			file.BookId, file.Chapter, ts.Verse, beginTS)
		outputPath := filepath.Join(tempDir, ts.AudioVerse)
		command = append(command, `-c`, `copy`, outputPath)
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
