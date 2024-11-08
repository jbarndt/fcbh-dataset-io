package tests

import (
	"bytes"
	"context"
	"database/sql"
	"dataset/controller"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

const (
	HOST       = `http://localhost:8080/`
	UPLOADHOST = `http://localhost:7777/upload`
	OUTPUT     = `/Users/gary/FCBH2024/systemtest/`
)

type SqliteTest struct {
	Query string
	Count int
}

// CLISqlTest requires a sqlite database as output to perform tests on the result
func DirectSqlTest(request string, tests []SqliteTest, t *testing.T) {
	database, status := controller.CLIProcessEntry([]byte(request))
	if status.IsErr {
		t.Fatal(status)
	}
	conn, err := sql.Open("sqlite3", database)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	var count int
	for _, tst := range tests {
		count = SelectScalarInt(conn, tst.Query, t)
		if count != tst.Count {
			t.Error("Count was " + strconv.Itoa(count) + ", expected " + strconv.Itoa(tst.Count) + " ON: " + tst.Query)
		}
	}
}

func SelectScalarInt(conn *sql.DB, sql string, t *testing.T) int {
	var count int
	rows, err := conn.Query(sql)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
	return count
}

func CLIExec(requestYaml string, t *testing.T) (string, string) {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), `request`+"_*.yaml")
	if err != nil {
		t.Error(err)
	}
	_, _ = file.Write([]byte(requestYaml))
	_ = file.Close()
	var cmd = exec.Command(`go`, `run`, `../controller/cli/dataset.go`, file.Name())
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

func ExtractFilename(yaml string) string {
	return ExtractFilenaame(yaml)
}

// Deprecated
func ExtractFilenaame(yaml string) string {
	start := strings.Index(yaml, `output_file:`) + 12
	end := strings.Index(yaml[start:], "\n")
	filename := strings.TrimSpace(yaml[start : start+end])
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
	records := strings.Split(string(content), "<tr>")
	return len(records) - 2
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
		t.Error(`AudioOTId expected`, audioOTId, `found`, ident.AudioOTId)
	}
	if ident.AudioNTId != audioNTId {
		t.Error(`AudioNTId expected`, audioNTId, `found`, ident.AudioNTId)
	}
	if ident.LanguageISO != language {
		t.Error(`LanguageISO expected`, language, `found`, ident.LanguageISO)
	}
}
