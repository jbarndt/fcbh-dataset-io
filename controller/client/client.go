package main

import (
	"bufio"
	"bytes"
	"context"
	"dataset/request"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	HOST = `http://localhost:8080/`
	//HOST       = `http:// 167.99.58.202:8080/`
	UPLOADHOST = `http://167.99.58.202:8080/upload`
)

type Arguments struct {
	yamlPath  string
	audioPath string
}

func main() {
	arg := GetArguments()
	yamlRequest, err := os.ReadFile(arg.yamlPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx := context.Background()
	reqDecoder := request.NewRequestDecoder(ctx)
	request, status := reqDecoder.Process(yamlRequest)
	if status.IsErr {
		fmt.Println(status)
		os.Exit(1)
	}
	var httpReq *http.Request
	httpReq = HttpPost(yamlRequest) // handle upload case ??
	statusCode := Response(request.OutputFile, httpReq)
	DisplayOutput(request.OutputFile)
	fmt.Println(statusCode)
}

func GetArguments() Arguments {
	var a Arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: client  filename.Yaml  [audioFile]")
		os.Exit(1)
	}
	a.yamlPath = os.Args[1]
	if len(os.Args) > 3 {
		a.audioPath = os.Args[2]
	}
	return a
}

func HttpPost(request []byte) *http.Request {
	//req, err := http.NewRequest("POST", HOST, bytes.NewReader(request))
	req, err := http.NewRequest("POST", HOST, bytes.NewBuffer(request))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-yaml")
	//req.Header.Set("Accept", "attachment; filename=\""+os.Args[1]+"\"")
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

func DisplayOutput(filename string) {
	if strings.HasSuffix(filename, ".sqlite") {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		fmt.Println(scanner.Text())
		if lineCount == 20 {
			break
		}
	}
}
