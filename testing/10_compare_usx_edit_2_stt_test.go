package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const CompareUsXTextEdit2STT = `is_new: no
dataset_name: AudioWhisperJson_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 10__compare_usx_edit_2_stt.html
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

func TestCompareUsXTextEdit2STTAPI(t *testing.T) {
	var bibleId = `ENGWEB`
	var req = strings.Replace(CompareUsXTextEdit2STT, `{bibleId}`, bibleId, 3)
	stdout, stderr := FCBHDatasetExec(req, t)
	fmt.Println(`STDOUT:`, stdout)
	fmt.Println(`STDERR:`, stderr)
	filename := ExtractFilenaame(req)
	fmt.Println("Filename", filename)
	count := NumHTMLFileLines(filename, t)
	expected := 22
	if count != expected {
		t.Error(`expected`, expected, `found`, count)
	}
}

func TestCompareUsXTextEdit2STT(t *testing.T) {
	var bibleId = `ENGWEB`
	ctx := context.Background()
	var req = strings.Replace(CompareUsXTextEdit2STT, `{bibleId}`, bibleId, 3)
	var control = controller.NewController(ctx, []byte(req))
	filename, status := control.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("Filename", filename)
	count := NumHTMLFileLines(filename, t)
	expected := 22
	if count != expected {
		t.Error(`expected`, expected, `found`, count)
	}
	identTest(`AudioWhisperJson_`+bibleId, t, request.TextSTT, ``,
		`ENGWEBN_TT`, ``, `ENGWEBN2DA`, `eng`)
}
