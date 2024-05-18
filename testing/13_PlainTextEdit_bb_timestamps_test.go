package testing

import (
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditBBTimestampsScript = `is_new: yes
dataset_name: PlainTextEditScript_{bibleId}
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
  bible_brain: yes
`

func TestPlainTextBBTimestampsScript(t *testing.T) {
	type try struct {
		bibleId   string
		textNtId  string
		audioNTId string
		language  string
		expected  int
	}
	var tests []try
	tests = append(tests, try{bibleId: "ENGWEB", expected: 8219, textNtId: "ENGWEBN_ET", audioNTId: "ENGWEBN2DA",
		language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	for _, tst := range tests {
		var req = strings.Replace(PlainTextEditBBTimestampsScript, `{bibleId}`, tst.bibleId, 2)
		var control = controller.NewController([]byte(req))
		filename, status := control.Process()
		if status.IsErr {
			t.Error(status)
		}
		fmt.Println(filename)
		numLines := NumCVSFileLines(filename, t)
		if numLines != tst.expected {
			t.Error(`Expected `, tst.expected, `records, got`, numLines)
		}
		identTest(`PlainTextEditScript_`+tst.bibleId, t, request.TextPlainEdit, ``,
			tst.textNtId, ``, tst.audioNTId, tst.language)
	}
}
