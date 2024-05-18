package testing

import (
	"bytes"
	"dataset/request"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const ClientTextTest = `is_new: yes
dataset_name: ClientUSXText_{bibleId}
bible_id: {bibleId}
text_data:
  bible_brain:
    text_usx_edit: yes
output_format:
  json: {json}
  csv: {csv}
`

func TestClientText(t *testing.T) {
	var bibleId = `ENGWEB`
	expected := 9588
	var req = strings.Replace(ClientTextTest, `{bibleId}`, bibleId, 2)
	req = strings.Replace(req, `{json}`, `yes`, 1)
	req = strings.Replace(req, `{csv}`, `no`, 1)
	stdout, stderr := ClientExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(stdout)
	numLines := NumJSONFileLines(filename, t)
	if numLines != expected {
		t.Error(`Expected `, expected, `records, got`, numLines)
	}
	identTest(`ClientUSXText_`+bibleId, t, request.TextUSXEdit, ``,
		`ENGWEBN_ET-usx`, ``, ``, `eng`)
}

func ClientExec(requestYaml string, t *testing.T) (string, string) {
	file, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), `request`+"_*.yaml")
	if err != nil {
		t.Error(err)
	}
	_, _ = file.Write([]byte(requestYaml))
	_ = file.Close()
	var cmd = exec.Command(`go`, `run`, `../controller/client/client.go`, file.Name())
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
