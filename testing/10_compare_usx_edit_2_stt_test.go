package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const CompareUsXTextEdit2STT = `is_new: no
dataset_name: AudioWhisperJson_{bibleId}
bible_id: {bibleId}
testament:
  nt_books: ['TIT']
compare:
  base_dataset: USX Text Edit Script_{bibleId}
  compare_settings: # Mark yes, all settings that apply
    lower_case: y
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

func TestCompareUsXTextEdit2STT(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(CompareUsXTextEdit2STT, `{bibleId}`, bibleId, 3)
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
