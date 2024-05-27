package testing

import (
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const PlainTextScript = `is_new: yes
dataset_name: PlainTextScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 01__plain_text_script.csv
text_data:
  bible_brain:
    text_plain: yes
`

func TestPlainTextScriptAPI(t *testing.T) {
	var cases []APITest
	cases = append(cases, APITest{BibleId: `ENGWEB`, Expected: 7959, Diff: 0})
	APITestUtility(PlainTextScript, cases, t)
}

func TestPlainTextScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var req = strings.Replace(PlainTextScript, `{bibleId}`, bibleId, 2)
	stdout, stderr := CLIExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilename(req)
	numLines := NumCVSFileLines(filename, t)
	count := 7959
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
	identTest(`PlainTextScript_`+bibleId, t, request.TextPlain, ``,
		`ENGWEBN_ET`, ``, ``, `eng`)
}
