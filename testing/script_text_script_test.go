package testing

import (
	"context"
	"dataset/controller"
	"dataset/db"
	"fmt"
	"strings"
	"testing"
)

const ScriptTextScript = `is_new: yes
request_name: ScriptTextScript
bible_id: {bibleId}
text_data:
  file: /Users/gary/FCBH2024/download/ATIWBT/ATIWBTN2ST.xlsx
output_format:
  sqlite: yes
`

func TestScriptTextScript(t *testing.T) {
	var bibleId = `ATIWBT`
	var request = strings.Replace(ScriptTextScript, `{bibleId}`, bibleId, 1)
	var control = controller.NewController([]byte(request))
	filename, status := control.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(filename)
	conn := db.NewDBAdapter(context.TODO(), filename)
	count, status := conn.CountScriptRows()
	if status.IsErr {
		t.Fatal(status)
	}
	if count != 10 {
		t.Error(`Expected `, count, `records, got`, count)
	}
}
