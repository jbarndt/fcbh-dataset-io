package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditScript = `is_new: yes
dataset_name: PlainTextEditScript_{bibleId}
bible_id: {bibleId}
text_data:
  bible_brain:
    text_plain_edit: yes
output_format:
  json: yes
`

func TestPlainTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(PlainTextEditScript, `{bibleId}`, bibleId, 2)
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

func TestPlainTextEditScript(t *testing.T) {
	var bibles = make(map[string]int)
	//bibles[`ENGWEB`] = 8250
	bibles[`ATIWBT`] = 8243
	for bibleId, expected := range bibles {
		var request = strings.Replace(PlainTextEditScript, `{bibleId}`, bibleId, 2)
		var control = controller.NewController([]byte(request))
		filename, status := control.Process()
		if status.IsErr {
			t.Error(status)
		}
		fmt.Println(filename)
		numLines := NumJSONFileLines(filename, t)
		if numLines != expected {
			t.Error(`Expected `, expected, `records, got`, numLines)
		}
	}
}
