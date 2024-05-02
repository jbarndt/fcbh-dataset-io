package testing

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const (
	HOST   = `http://localhost:8080/`
	OUTPUT = `/Users/gary/FCBH2024/systemtest/`
)

func HttpPost(request string, name string, t *testing.T) ([]byte, int) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", HOST, bytes.NewReader([]byte(request)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-yaml")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	// can close go here
	filePath := filepath.Join(OUTPUT, name)
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.Write(body)
	if err != nil {
		t.Fatal(err)
	}
	_ = file.Close()
	return body, resp.StatusCode
}

func NumCVSLines(content []byte, t *testing.T) int {
	memFile := bytes.NewReader(content)
	reader := csv.NewReader(memFile)
	count := 0
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		count++
	}
	return count
}
