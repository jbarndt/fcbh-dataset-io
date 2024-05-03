package testing

import (
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditScript = `IsNew: yes
RequestName: PlainTextEditScript
BibleId: {bibleId}
TextData:
  BibleBrain:
    TextPlainEdit: yes
OutputFormat:
  JSON: yes
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
