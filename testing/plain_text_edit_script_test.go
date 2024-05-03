package testing

import (
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditScript = `Required:
  IsNew: yes
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
	start := strings.Index(stdout, `Success: `) + 9
	end := strings.Index(stdout[start:], "\n")
	filename := stdout[start : start+end]
	fmt.Println(filename)
	numLines := NumJSONFileLines(filename, t)
	count := 8250
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
}
