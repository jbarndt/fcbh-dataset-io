package mms

import (
	"bufio"
	"context"
	"dataset"
	log "dataset/logger"
	"github.com/garygriswold/lang_tree/search"
	"io"
	"os/exec"
	"strings"
)

// Check that language is supported by mms_asr, and return alternate if it is not
func checkLanguage(ctx context.Context, lang string, sttLang string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	if sttLang != `` {
		result = sttLang
	} else {
		var tree = search.NewLanguageTree(ctx)
		err := tree.Load()
		if err != nil {
			status = log.Error(ctx, 500, err, `Error loading language`)
			return result, status
		}
		langs, distance, err2 := tree.Search(strings.ToLower(lang), "mms_asr")
		if err2 != nil {
			status = log.Error(ctx, 500, err2, `Error Searching for language`)
		}
		if len(langs) > 0 {
			result = langs[0]
			log.Info(ctx, `Using language`, result, "distance:", distance)
		} else {
			status = log.ErrorNoErr(ctx, 400, `No compatible language code was found for`, lang)
		}
	}
	return result, status
}

// callPythonScript will exec the python script, and setup pipes on stdin and stdout
func callStdIOScript(ctx context.Context, command string, arg ...string) (io.Writer, io.Reader, dataset.Status) {
	var writer io.Writer
	var reader io.Reader
	var status dataset.Status
	cmd := exec.Command(command, arg...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to open stdin for writing to Fasttext`)
		return writer, reader, status
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to open stdout for writing to Fasttext`)
		return writer, reader, status
	}
	err = cmd.Start()
	if err != nil {
		status = log.Error(ctx, 500, err, `Unable to start writing to Fasttext`)
		return writer, reader, status
	}
	writer = bufio.NewWriterSize(stdin, 4096)
	reader = bufio.NewReaderSize(stdout, 4096)
	return writer, reader, status
}
