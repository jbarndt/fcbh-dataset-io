package controller

import (
	"context"
	"dataset"
	"io"
	"os"
	"strings"
)

func CLIProcessEntry(yaml []byte) (string, dataset.Status) {
	var result string
	var ctx = context.WithValue(context.Background(), `runType`, `cli`)
	var control = NewController(ctx, yaml)
	filename, status := control.Process()
	if status.IsErr {
		result = status.String()
	}
	outputFile := findOutputFilename(yaml)
	var err error
	if strings.HasSuffix(outputFile, ".sqlite") {
		err = copyFile(filename, outputFile)
	} else {
		err = os.Rename(filename, outputFile)
	}
	if err != nil {
		result = filename
	} else {
		result = outputFile
	}
	return result, status
}

func findOutputFilename(request []byte) string {
	var result string
	req := string(request)
	start := strings.Index(req, `output_file:`) + 12
	end := strings.Index(req[start:], "\n")
	result = strings.TrimSpace(req[start : start+end])
	return result
}

func copyFile(src, dst string) error {
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	return err
}
