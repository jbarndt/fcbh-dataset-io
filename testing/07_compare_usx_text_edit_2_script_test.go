package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const CompareUsXTextEdit2Script = `is_new: no
dataset_name: ScriptTextScript_{bibleId}
bible_id: {bibleId}
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
	var request = strings.Replace(CompareUsXTextEdit2Script, `{bibleId}`, bibleId, 3)
	var control = controller.NewController([]byte(request))
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
}
