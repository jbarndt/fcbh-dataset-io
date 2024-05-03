package main

import (
	"dataset/controller"
	"fmt"
	"os"
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
	fmt.Fprintln(os.Stdout, `Success:`, filename)
}
