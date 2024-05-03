package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const USXTextEditScript = `Required:
  IsNew: yes
  RequestName: USX Text Edit Script
  BibleId: {bibleId}
TextData:
  BibleBrain:
    TextUSXEdit: yes
OutputFormat:
  JSON: yes
`

func TestUSXTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(USXTextEditScript, `{bibleId}`, bibleId, 1)
	stdout, stderr := CLIExec(request, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(stdout)
	numLines := NumJSONFileLines(filename, t)
	count := 8250
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
