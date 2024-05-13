package testing

import (
	"bytes"
	"context"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	HOST       = `http://localhost:8080/`
	UPLOADHOST = `http://localhost:8080/upload`
	OUTPUT     = `/Users/gary/FCBH2024/systemtest/`
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

func CLIExec(requestYaml string, t *testing.T) (string, string) {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), `request`+"_*.yaml")
	if err != nil {
		t.Error(err)
	}
	_, _ = file.Write([]byte(requestYaml))
	_ = file.Close()
	var cmd = exec.Command(`go`, `run`, `../controller/client/dataset.go`, file.Name())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		t.Error(err.Error())
	}
	_ = os.Remove(file.Name())
	return stdoutBuf.String(), stderrBuf.String()
}

func ExtractFilenaame(stdout string) string {
	start := strings.Index(stdout, `Success: `) + 9
	end := strings.Index(stdout[start:], "\n")
	filename := stdout[start : start+end]
	return filename
}

func NumCVSLines(content []byte, t *testing.T) int {
	file := bytes.NewReader(content)
	reader := csv.NewReader(file)
	return numCVSLineGeneric(reader, t)
}

func NumCVSFileLines(filename string, t *testing.T) int {
	file, err := os.Open(filename)
	if err != nil {
		t.Error(err)
	}
	reader := csv.NewReader(file)
	return numCVSLineGeneric(reader, t)
}

func numCVSLineGeneric(reader *csv.Reader, t *testing.T) int {
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

func NumJSONFileLines(filename string, t *testing.T) int {
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	return NumJSONLines(content, t)
}

func NumJSONLines(content []byte, t *testing.T) int {
	var response []map[string]any
	err := json.Unmarshal(content, &response)
	if err != nil {
		t.Error(err)
	}
	count := len(response)
	return count
}

func NumHTMLFileLines(filename string, t *testing.T) int {
	//count := 0
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	records := strings.Split(string(content), "\n")
	return len(records)
}

func identTest(name string, t *testing.T, textType request.MediaType, textOTId string,
	textNTId string, audioOTId string, audioNTId string, language string) {
	user, _ := fetch.GetTestUser()
	conn, status := db.NewerDBAdapter(context.TODO(), false, user.Username, name)
	if status.IsErr {
		t.Fatal(status)
	}
	ident, status := conn.SelectIdent()
	if status.IsErr {
		t.Fatal(status)
	}
	conn.Close()
	if ident.TextSource != textType {
		t.Error(`TextSource expected`, textType, `found`, ident.TextSource)
	}
	if ident.TextOTId != textOTId {
		t.Error(`TextOTId expected`, textOTId, `found`, ident.TextOTId)
	}
	if ident.TextNTId != textNTId {
		t.Error(`TextNTId expected`, textNTId, `found`, ident.TextNTId)
	}
	if ident.AudioOTId != audioOTId {
		t.Error(`AudioOTId expected`, textOTId, `found`, ident.AudioOTId)
	}
	if ident.AudioNTId != audioNTId {
		t.Error(`AudioNTId expected`, textNTId, `found`, ident.AudioNTId)
	}
	if ident.LanguageISO != language {
		t.Error(`LanguageISO expected`, language, `found`, ident.LanguageISO)
	}
}
