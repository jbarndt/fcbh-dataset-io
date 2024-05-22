package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const USXTextEditScript = `is_new: yes
dataset_name: USX Text Edit Script_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 03__usx_text_edit_script.json
text_data:
  bible_brain:
    text_usx_edit: yes
`

func TestUSXTextEditScriptCLI(t *testing.T) {
	var bibleId = `ENGWEB`
	var req = strings.Replace(USXTextEditScript, `{bibleId}`, bibleId, 2)
	stdout, stderr := CLIExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(req)
	numLines := NumJSONFileLines(filename, t)
	expected := 9588
	if numLines != expected {
		t.Error(`Expected `, expected, `records, got`, numLines)
	}
	identTest(`USX_Text_Edit_Script_`+bibleId, t, request.TextUSXEdit, ``,
		`ENGWEBN_ET-usx`, ``, ``, `eng`)
}

func TestUSXTextEditScript(t *testing.T) {
	var bibleId = `ATIWBT`
	ctx := context.Background()
	var req = strings.Replace(USXTextEditScript, `{bibleId}`, bibleId, 2)
	var control = controller.NewController(ctx, []byte(req))
	filename, status := control.Process()
	if status.IsErr {
		t.Error(status)
	}
	fmt.Println(filename)
	numLines := NumJSONFileLines(filename, t)
	count := 8254
	if numLines != count {
		t.Error(`Expected `, count, `records, got`, numLines)
	}
	identTest(`USX_Text_Edit_Script_`+bibleId, t, request.TextUSXEdit, ``,
		`ATIWBTN_ET-usx`, ``, ``, `ati`)
}
