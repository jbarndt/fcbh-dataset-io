package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const CompareUsXTextEdit2Script = `is_new: no
dataset_name: ScriptTextScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 07__compare_usx_text_edit_2_script.html
compare:
  base_dataset: USX Text Edit Script_{bibleId}
  compare_settings: # Mark yes, all settings that apply
    lower_case: n
    remove_prompt_chars: y
    remove_punctuation: y
    double_quotes: 
      remove: y
      normalize:
    apostrophe: 
      remove: y
      normalize:
    hyphen:
      remove: y
      normalize:
    diacritical_marks: 
      remove:
      normalize_nfc: 
      normalize_nfd: y
      normalize_nfkc:
      normalize_nfkd:
`

func TestCompareUsXTextEdit2Script(t *testing.T) {
	var bibleId = `ATIWBT`
	ctx := context.Background()
	var req = strings.Replace(CompareUsXTextEdit2Script, `{bibleId}`, bibleId, 3)
	var control = controller.NewController(ctx, []byte(req))
	filename, status := control.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("Filename", filename)
	count := NumHTMLFileLines(filename, t)
	expected := 18
	if count != expected {
		t.Error(`expected`, expected, `found`, count)
	}
	identTest(`ScriptTextScript_`+bibleId, t, request.TextScript, ``,
		`ATIWBTN2ST`, ``, ``, `ati`)
}
