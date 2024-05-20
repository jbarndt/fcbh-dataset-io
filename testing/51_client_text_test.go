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
username: GaryNTest
email: gary@shortsands.com
output_file: 51__client_text.json
text_data:
  bible_brain:
    text_usx_edit: yes
`

func TestClientText(t *testing.T) {
	var bibleId = `ENGWEB`
	expected := 9588
	var req = strings.Replace(ClientTextTest, `{bibleId}`, bibleId, 2)
	stdout, stderr := ClientExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(req)
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
