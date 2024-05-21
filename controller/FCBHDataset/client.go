package main

import (
	"bufio"
	"bytes"
	"context"
	"dataset/request"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func main() {
	cfg := GetConfig()
	yamlPath := GetArguments()
	yamlRequest, err := os.ReadFile(yamlPath)
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
	if request.AudioData.POST != `` {
		httpReq = HttpMultiPost(cfg, yamlRequest, request.AudioData.POST, "audio")
	} else if request.TextData.POST != `` {
		httpReq = HttpMultiPost(cfg, yamlRequest, request.TextData.POST, "text")
	} else {
		httpReq = HttpPost(cfg, yamlRequest)
	}
	statusCode := Response(request.OutputFile, httpReq)
	DisplayOutput(request.OutputFile)
	fmt.Println(statusCode, request.OutputFile)
}

func GetArguments() string {
	if len(os.Args) < 2 {
		fmt.Println("Usage: FCBHDataset  filename.yaml")
		os.Exit(1)
	}
	return os.Args[1]
}

func HttpPost(cfg Config, request []byte) *http.Request {
	req, err := http.NewRequest("POST", cfg.Host, bytes.NewBuffer(request))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-yaml")
	return req
}

func HttpMultiPost(cfg Config, yamlRequest []byte, filePath string, fType string) *http.Request {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	filePart, err := writer.CreateFormFile(fType, filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	yamlPart, err := writer.CreateFormFile("yaml", "request.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(yamlPart, bytes.NewBuffer(yamlRequest))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_ = writer.Close()
	req, err := http.NewRequest("POST", cfg.Host+`/upload`, body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
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

func Catch(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
