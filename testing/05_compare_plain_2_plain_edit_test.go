package testing

import (
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const ComparePlain2PlainEditScript = `is_new: no
dataset_name: PlainTextScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 05__compare_plain_2_edit.html
compare:
  base_dataset: PlainTextEditScript_{bibleId}
  compare_settings: # Mark yes, all settings that apply
    lower_case: 
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

func TestComparePlain2PlainEditScript(t *testing.T) {
	var bibleId = `ENGWEB`
	var req = strings.Replace(ComparePlain2PlainEditScript, `{bibleId}`, bibleId, 3)
	var control = controller.NewController([]byte(req))
	filename, status := control.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("Filename", filename)
	count := NumHTMLFileLines(filename, t)
	expected := 276
	if count != expected {
		t.Error(`expected`, expected, `found`, count)
	}
	identTest(`PlainTextScript_`+bibleId, t, request.TextPlain, ``,
		`ENGWEBN_ET`, ``, ``, `eng`)
}
