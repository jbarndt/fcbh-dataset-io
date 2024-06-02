package testing

import (
	"dataset/request"
	"testing"
)

const PlainTextEditTSBBScript = `is_new: yes
dataset_name: PlainTextEditTSBBScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 14__plain_text_edit_bb.csv
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps: 
  bible_brain: yes
`

func TestPlainTextBBTimestampsScript(t *testing.T) {
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 8219, TextNtId: "ENGWEBN_ET",
		TextType: request.TextPlainEdit, AudioNTId: "ENGWEBN2DA", Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(PlainTextEditTSBBScript, tests, t)
}

// ENGWEB BB timestamps
// select avg(script_end_ts-script_begin_ts) from scripts where script_end_ts != 0.0
// = 8.37511692230324
