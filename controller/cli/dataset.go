package main

import (
	"context"
	"dataset"
	"dataset/controller"
	log "dataset/logger"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stdout, "Usage: dataset  request.yaml")
		os.Exit(1)
	}
	outputFile, status := MainProcess(os.Args[1])
	if status.IsErr {
		_, _ = fmt.Fprintln(os.Stderr, status.String())
		os.Exit(1)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, `Success:`, outputFile)
	}
}

func MainProcess(yamlPath string) (string, dataset.Status) {
	var result string
	var content, err = os.ReadFile(yamlPath)
	if err != nil {
		return result, log.Error(context.Background(), 400, err, `Error reading yaml request file.`)
	}
	var ctx = context.WithValue(context.Background(), `runType`, `cli`)
	var control = controller.NewController(ctx, content)
	filename, status := control.Process()
	if status.IsErr {
		result = status.String()
	}
	outputFile := findOutputFilename(content)
	err = os.Rename(filename, outputFile)
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
