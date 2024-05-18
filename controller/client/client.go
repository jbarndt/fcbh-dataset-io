package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	HOST = `http://localhost:8080/`
	//HOST       = `http:// 167.99.58.202:8080/`
	UPLOADHOST = `http://167.99.58.202:8080/upload`
)

type Arguments struct {
	yamlPath  string
	audioPath string
	outPath   string
}

func main() {
	arg := GetArguments()
	// read yaml file
	yamlFile, err := os.ReadFile(arg.yamlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var req *http.Request
	if arg.audioPath == `` {
		req = HttpPost(yamlFile)
	} else {

	}
	// get arguments
	// decide if typical or post audio
	// call correct http routine
	// receive
	status := Response(arg.outPath, req)
	fmt.Println(status)
}

func GetArguments() Arguments {
	var a Arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: client  filename.Yaml  [audioFile]  outputFile")
		os.Exit(1)
	}
	a.yamlPath = os.Args[1]
	if len(os.Args) > 3 {
		a.audioPath = os.Args[2]
		a.outPath = os.Args[3]
	} else {
		a.outPath = os.Args[2]
	}
	return a
}

// read the command line
// first parameter is yaml file
// next paramter can be audio file
// last parameter is outout filepath

func HttpPost(request []byte) *http.Request {

	req, err := http.NewRequest("POST", HOST, bytes.NewReader(request))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//return body, statusCode
	}
	req.Header.Set("Content-Type", "application/x-yaml")
	return req
}

func HttpPostFile(request string, audioPath string) *http.Request {
	req := http.Request{}
	return &req
}

func Response(outPath string, req *http.Request) string {
	var status string
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		//return body, statusCode
	}
	status = resp.Status
	defer resp.Body.Close()
	file, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
		return status
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return status
	}
	_ = file.Close()
	return status
}
