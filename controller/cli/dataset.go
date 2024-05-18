package main

import (
	"dataset/controller"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: dataset  request.yaml")
		os.Exit(1)
	}
	var yamlPath = os.Args[1]
	var content, err = os.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var control = controller.NewController(content)
	filename, status := control.Process()
	if status.IsErr {
		fmt.Fprintln(os.Stderr, status.String())
		fmt.Fprintln(os.Stderr, `Error File:`, filename)
		os.Exit(1)
	}
	//fmt.Fprintln(os.Stdout, `Success:`, filename)
	outputFile := findOutputFilename(content)
	err = os.Rename(filename, outputFile)
	if err != nil {
		fmt.Fprintln(os.Stdout, `Success:`, filename)
	} else {
		fmt.Fprintln(os.Stdout, `Success:`, outputFile)
	}
}

func findOutputFilename(request []byte) string {
	var result string
	req := string(request)
	start := strings.Index(req, `output_file:`) + 12
	end := strings.Index(req[start:], "\n")
	result = req[start : start+end]
	return result
}
