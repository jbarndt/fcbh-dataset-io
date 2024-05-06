package testing

import (
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditScript = `is_new: yes
dataset_name: PlainTextEditScript
bible_id: {bibleId}
text_data:
  bible_brain:
    text_plain_edit: yes
output_format:
  json: yes
`

func TestPlainTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(PlainTextEditScript, `{bibleId}`, bibleId, 1)
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
