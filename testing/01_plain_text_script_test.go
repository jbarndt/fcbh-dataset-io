package testing

import (
	"fmt"
	"strings"
	"testing"
)

const PlainTextScript = `is_new: yes
dataset_name: PlainTextScript_{bibleId}
bible_id: {bibleId}
text_data:
  bible_brain:
    text_plain: yes
output_format:
  csv: yes
`

func TestPlainTextScriptAPI(t *testing.T) {
	var cases = make(map[string]int)
	cases[`ENGWEB`] = 7959
	for bibleId, count := range cases {
		var request = strings.Replace(PlainTextScript, `{bibleId}`, bibleId, 2)
		csvResp, statusCode := HttpPost(request, `PlainTextScript.csv`, t)
		fmt.Printf("Response status: %d\n", statusCode)
		//fmt.Println("Response body:", string(csvResp))
		numLines := NumCVSLines(csvResp, t)
		if numLines != count {
			t.Error(`Expected `, count, `records, got`, numLines)
		}
	}
}

func TestPlainTextScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(PlainTextScript, `{bibleId}`, bibleId, 2)
	stdout, stderr := CLIExec(request, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(stdout)
	numLines := NumCVSFileLines(filename, t)
	count := 7959
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
}
