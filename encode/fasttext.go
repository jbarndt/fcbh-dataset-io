package encode

import (
	"bufio"
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/**
FastText
https://github.com/facebookresearch/fastText?tab=readme-ov-file
*/

type FastText struct {
	ctx       context.Context
	conn      db.DBAdapter
	bibleId   string // needed
	filesetId string // needed
}

func NewFastText(ctx context.Context, conn db.DBAdapter) FastText {
	var f FastText
	f.ctx = ctx
	f.conn = conn
	return f
}

func (f *FastText) Process() dataset.Status {
	var words, status = f.conn.SelectWords()
	if status.IsErr {
		return status
	}
	inputFile, status := f.createFile(words)
	if status.IsErr {
		return status
	}
	outputModel, status := f.executeFastText(inputFile)
	if status.IsErr {
		return status
	}
	words, status = f.getWordEncodings(outputModel, words)
	if status.IsErr {
		return status
	}
	f.conn.UpdateWordEncodings(words)
	return status
}

func (f *FastText) createFile(words []db.Word) (string, dataset.Status) {
	var status dataset.Status
	var fp, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), `fasttextinput`)
	if err != nil {
		status = log.Error(f.ctx, 500, err, `Unable to open temp file for fasttext`)
		return ``, status
	}
	for _, word := range words {
		_, err = fp.WriteString(word.Word)
		if err != nil {
			status = log.Error(f.ctx, 500, err, `Error while writing to fasttext input file`)
			return fp.Name(), status
		}
	}
	fp.Close()
	return fp.Name(), status
}

func (f *FastText) executeFastText(inputFile string) (string, dataset.Status) {
	var status dataset.Status
	fastTextExe := os.Getenv("FCBH_FASTTEXT_EXE")
	model := `skipgram` // or `cbow
	outputModel := strings.Replace(f.conn.DatabasePath, `.db`, `_fasttext`, 1)
	cmd := exec.Command(fastTextExe, model, `-input`, inputFile, `-output`, outputModel)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		status = log.Error(f.ctx, 500, err, `Error executing Fasttext`)
	}
	if stderrBuf.Len() > 0 {
		fmt.Println("STDERR", stderrBuf.String())
	}
	if stdoutBuf.Len() > 0 {
		fmt.Println("STDOUT", stdoutBuf.String())
	}
	return outputModel, status
}

func (f *FastText) getWordEncodings(model string, words []db.Word) ([]db.Word, dataset.Status) {
	var status dataset.Status
	fastTextExe := os.Getenv("FCBH_FASTTEXT_EXE")
	cmd := exec.Command(fastTextExe, `print-word-vectors`, model+`.bin`)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		status = log.Error(f.ctx, 500, err, `Unable to open stdin for writing to Fasttext`)
		return words, status
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		status = log.Error(f.ctx, 500, err, `Unable to open stdout for writing to Fasttext`)
		return words, status
	}
	err = cmd.Start()
	if err != nil {
		status = log.Error(f.ctx, 500, err, `Unable to start writing to Fasttext`)
		return words, status
	}
	reader := bufio.NewReader(stdout)
	for i, word := range words {
		if word.TType == `W` {
			_, err := io.WriteString(stdin, word.Word+"\n")
			if err != nil {
				status = log.Error(f.ctx, 500, err, `Error writing to Fasttext model`)
				return words, status
			}
			line, err := reader.ReadString('\n')
			if err != nil {
				status = log.Error(f.ctx, 500, err, `Error reading from Fasttext model`)
				return words, status
			}
			parts := strings.Split(strings.TrimSpace(line), ` `)
			var encoding = make([]float64, 0, len(parts))
			for _, strNum := range parts[1:] {
				num, err := strconv.ParseFloat(strNum, 64)
				if err != nil {
					status = log.Error(f.ctx, 500, err, `Error converting encoding to float`)
					return words, status
				}
				encoding = append(encoding, num)
			}
			word.WordEncoded = encoding
		}
		words[i] = word
	}
	return words, status
}
