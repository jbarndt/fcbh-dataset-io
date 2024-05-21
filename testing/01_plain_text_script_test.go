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
	var cases = make(map[string]int)
	cases[`ENGWEB`] = 7959
	for bibleId, count := range cases {
		var req = strings.Replace(PlainTextScript, `{bibleId}`, bibleId, 2)
		csvResp, statusCode := HttpPost(req, `PlainTextScript.csv`, t)
		fmt.Printf("Response status: %d\n", statusCode)
		//fmt.Println("Response body:", string(csvResp))
		numLines := NumCVSLines(csvResp, t)
		if numLines != count {
			t.Error(`Expected `, count, `records, got`, numLines)
		}
		identTest(`PlainTextScript_`+bibleId, t, request.TextPlain, ``,
			`ENGWEBN_ET`, ``, ``, `eng`)
	}
}

func TestPlainTextScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var req = strings.Replace(PlainTextScript, `{bibleId}`, bibleId, 2)
	stdout, stderr := CLIExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(req)
	numLines := NumCVSFileLines(filename, t)
	count := 7959
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
	identTest(`PlainTextScript_`+bibleId, t, request.TextPlain, ``,
		`ENGWEBN_ET`, ``, ``, `eng`)
}
