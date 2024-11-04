package controller

import (
	"context"
	"dataset"
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
	err := os.Rename(filename, outputFile)
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
