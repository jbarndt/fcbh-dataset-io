package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const ComparePlain2PlainEditScript = `is_new: no
dataset_name: PlainTextScript
bible_id: {bibleId}
compare:
  base_dataset: PlainTextEditScript
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

//output_format:
//  html: yes

func TestComparePlain2PlainEditScript(t *testing.T) {
	var bibleId = `ENGWEB`
	var request = strings.Replace(ComparePlain2PlainEditScript, `{bibleId}`, bibleId, 1)
	var control = controller.NewController([]byte(request))
	filename, status := control.Process()
	fmt.Println("Filename", filename)
	if status.IsErr {
		t.Fatal(status)
	}

	//numLines := NumJSONFileLines(filename, t)
	//count := 9568
	//if numLines != count {
	//	t.Error(`Expected `, count, `records, got`, numLines)
	//}
}
