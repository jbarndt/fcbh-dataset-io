package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const USXTextEditScript = `is_new: yes
request_name: USX Text Edit Script
bible_id: {bibleId}
text_data:
  bible_brain:
    text_usx_edit: yes
output_format:
  json: yes
`

func TestUSXTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(USXTextEditScript, `{bibleId}`, bibleId, 1)
	stdout, stderr := CLIExec(request, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(stdout)
	numLines := NumJSONFileLines(filename, t)
	count := 9568
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
}

func TestUSXTextEditScript(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(USXTextEditScript, `{bibleId}`, bibleId, 1)
	var control = controller.NewController([]byte(request))
	filename, status := control.Process()
	if status.IsErr {
		t.Error(status)
	}
	fmt.Println(filename)
	numLines := NumJSONFileLines(filename, t)
	count := 9568
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
}
