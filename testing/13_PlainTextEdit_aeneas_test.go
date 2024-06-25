package testing

import (
	"dataset/request"
	"testing"
)

const PlainTextEditBBTimestampsScript = `is_new: yes
dataset_name: PlainTextEditTSScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 13__plain_text_edit_bb_timestamps.csv
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  aeneas: yes
testament:
  nt_books: ['1JN']
`

func TestPlainTextAeneasTimestampsScriptAPI(t *testing.T) {
	var tests []APITest
	tests = append(tests, APITest{BibleId: `ENGWEB`, Expected: 111, Diff: 0})
	tests = append(tests, APITest{BibleId: `ATIWBT`, Expected: 111, Diff: 0})
	APITestUtility(PlainTextEditBBTimestampsScript, tests, t)
}
func TestPlainTextAeneasTimestampsScript(t *testing.T) {
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 111, TextNtId: "ENGWEBN_ET",
		TextType: request.TextPlainEdit, AudioNTId: "ENGWEBN2DA", Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(PlainTextEditBBTimestampsScript, tests, t)
}
