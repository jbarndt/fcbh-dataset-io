package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditScript = `is_new: no
dataset_name: PlainTextEditScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 02__plain_text_edit_script.json
text_data:
  bible_brain:
    text_plain_edit: yes
`

func TestPlainTextEditScriptAPI(t *testing.T) {
	var cases []APITest
	cases = append(cases, APITest{BibleId: `ENGWEB`, Expected: 8218})
	cases = append(cases, APITest{BibleId: `ATIWBT`, Expected: 8216})
	APITestUtility(PlainTextEditScript, cases, t)
}

func TestPlainTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var expected = 8218
	var req = strings.Replace(PlainTextEditScript, `{bibleId}`, bibleId, 2)
	stdout, stderr := CLIExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilename(req)
	numLines := NumJSONFileLines(filename, t)
	if numLines != expected {
		t.Error(`Expected `, expected, `records, got`, numLines)
	}
	identTest(`PlainTextEditScript_`+bibleId, t, request.TextPlainEdit, ``,
		`ENGWEBN_ET`, ``, ``, `eng`)
}

func TestPlainTextEditScript(t *testing.T) {
	type test struct {
		bibleId  string
		expected int
		textNtId string
		language string
	}
	var tests []test
	tests = append(tests, test{bibleId: "ENGWEB", expected: 8218, textNtId: "ENGWEBN_ET", language: "eng"})
	tests = append(tests, test{bibleId: "ATIWBT", expected: 8216, textNtId: "ATIWBTN_ET", language: "ati"})
	ctx := context.Background()
	for _, tst := range tests {
		var req = strings.Replace(PlainTextEditScript, `{bibleId}`, tst.bibleId, 2)
		var control = controller.NewController(ctx, []byte(req))
		filename, status := control.Process()
		if status.IsErr {
			t.Error(status)
		}
		fmt.Println(filename)
		numLines := NumJSONFileLines(filename, t)
		if numLines != tst.expected {
			t.Error(`Expected `, tst.expected, `records, got`, numLines)
		}
		identTest(`PlainTextEditScript_`+tst.bibleId, t, request.TextPlainEdit, ``,
			tst.textNtId, ``, ``, tst.language)
	}
}
