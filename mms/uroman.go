package mms

import (
	"bytes"
	"context"
	"dataset"
	log "dataset/logger"
	"os"
	"os/exec"
	"strings"
)

// uroman.go requires a pip install uroman, but it uses the uroman.pl that is included
// https://github.com/isi-nlp/uroman/tree/master

func URoman(ctx context.Context, lang string, text []string) ([]string, dataset.Status) {
	var result []string
	var status dataset.Status
	uromanPath := os.Getenv(`FCBH_UROMAN_EXE`)
	cmd := exec.Command(uromanPath, "-l", lang)
	inputStr := strings.Join(text, "\n")
	cmd.Stdin = strings.NewReader(inputStr)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(ctx, 500, err, "Command to execute uroman.pl failed.", stderrBuf.String())
		return result, status
	}
	result = strings.Split(strings.TrimSpace(stdoutBuf.String()), "\n")
	return result, status
}
