package main

import (
	"dataset/controller"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stdout, "Usage: dataset  request.yaml")
		os.Exit(1)
	}
	var content, err = os.ReadFile(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stdout, "Error reading yaml request file.")
		os.Exit(1)
	}
	outputFile, status := controller.CLIProcessEntry(content)
	if status.IsErr {
		_, _ = fmt.Fprintln(os.Stderr, status.String())
		os.Exit(1)
	} else {
		_, _ = fmt.Fprintln(os.Stdout, `Success:`, outputFile)
	}
}
