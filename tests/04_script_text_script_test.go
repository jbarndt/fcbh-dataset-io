package tests

import (
	"context"
	"dataset/controller"
	"dataset/db"
	"dataset/decode_yaml/request"
	"fmt"
	"strings"
	"testing"
)

const ScriptTextScript = `is_new: yes
dataset_name: ScriptTextScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output:
  sqlite: yes
text_data:
  file: /Users/gary/FCBH2024/download/ATIWBT/ATIWBTN2ST.xlsx
`

func TestScriptTextScriptAPI(t *testing.T) {
	var cases []APITest
	cases = append(cases, APITest{BibleId: `ATIWBT`, Expected: 9747})
	APITestUtility(ScriptTextScript, cases, t)
}

func TestScriptTextScript(t *testing.T) {
	var bibleId = `ATIWBT`
	ctx := context.Background()
	var req = strings.Replace(ScriptTextScript, `{bibleId}`, bibleId, 2)
	var control = controller.NewController(ctx, []byte(req))
	filename, status := control.Process()
	if status != nil {
		t.Fatal(status)
	}
	fmt.Println("Filename:", filename)
	conn := db.NewDBAdapter(context.TODO(), filename)
	count, status := conn.CountScriptRows()
	if status != nil {
		t.Fatal(status)
	}
	var expected = 9747
	if count != expected {
		t.Error(`Expected `, expected, `records, got`, count)
	}
	identTest(`ScriptTextScript_`+bibleId, t, request.TextScript, ``,
		`ATIWBTN2ST`, ``, ``, `ati`)
}
