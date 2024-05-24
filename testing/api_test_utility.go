package testing

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type APITest struct {
	BibleId  string
	Expected int
	Diff     int
}

func APITestUtility(request string, tests []APITest, t *testing.T) {
	for _, tst := range tests {
		var req = strings.Replace(request, `{bibleId}`, tst.BibleId, 3)
		stdout, stderr := FCBHDatasetExec(req, t)
		fmt.Println(`STDOUT:`, stdout)
		fmt.Println(`STDERR:`, stderr)
		filename := ExtractFilename(req)
		numLines := NumFileLines(filename, t)
		if numLines >= 0 {
			if numLines < (tst.Expected-tst.Diff) || numLines > (tst.Expected+tst.Diff) {
				t.Error(`Expected `, tst.Expected, `records, got`, numLines)
			}
		}
	}
}

func FCBHDatasetExec(requestYaml string, t *testing.T) (string, string) {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), `request`+"_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.Write([]byte(requestYaml))
	_ = file.Close()
	var cmd = exec.Command(`go`, `run`, `../controller/FCBHDataset`, file.Name())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		fmt.Println(stderrBuf.String())
		t.Fatal(err.Error())
	}
	_ = os.Remove(file.Name())
	return stdoutBuf.String(), stderrBuf.String()
}

func NumFileLines(filename string, t *testing.T) int {
	extension := filepath.Ext(filename)
	var numLines int
	switch extension {
	case ".json":
		numLines = NumJSONFileLines(filename, t)
	case ".csv":
		numLines = NumCVSFileLines(filename, t)
	case ".html":
		numLines = NumHTMLFileLines(filename, t)
	case ".sqlite":
		numLines = -1
	default:
		t.Fatal("Unexpected output_file type:", extension)
	}
	return numLines
}
