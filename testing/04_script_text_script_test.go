package testing

import (
	"context"
	"dataset/controller"
	"dataset/db"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const ScriptTextScript = `is_new: yes
dataset_name: ScriptTextScript_{bibleId}
bible_id: {bibleId}
text_data:
  file: /Users/gary/FCBH2024/download/ATIWBT/ATIWBTN2ST.xlsx
output_format:
  sqlite: yes
`

func TestScriptTextScript(t *testing.T) {
	var bibleId = `ATIWBT`
	var req = strings.Replace(ScriptTextScript, `{bibleId}`, bibleId, 2)
	var control = controller.NewController([]byte(req))
	filename, status := control.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("Filename:", filename)
	conn := db.NewDBAdapter(context.TODO(), filename)
	count, status := conn.CountScriptRows()
	if status.IsErr {
		t.Fatal(status)
	}
	var expected = 9747
	if count != expected {
		t.Error(`Expected `, expected, `records, got`, count)
	}
	identTest(`ScriptTextScript_`+bibleId, t, request.TextScript, ``,
		`ATIWBTN2ST`, ``, ``, `ati`)
}
