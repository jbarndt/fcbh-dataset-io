package testing

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

//const (
//	HOST = "http://127.0.0.1:8080"
//)

func SubmitRequestPost(filePath string, yamlPath string, t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	audioFile, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer audioFile.Close()
	audioPart, err := writer.CreateFormFile("audio", filePath)
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(audioPart, audioFile)
	yamlFile, err := os.Open(yamlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer yamlFile.Close()
	yamlPart, err := writer.CreateFormFile("metadata", "request.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(yamlPart, yamlFile)
	if err != nil {
		t.Fatal(err)
	}
	_ = writer.Close()
	request, err := http.NewRequest("POST", HOST, body)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "application/json")

	// Execute the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Status Code:", response.StatusCode)
	fmt.Println("Response Body:", responseBody.String())
}
