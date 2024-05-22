package testing

import (
	"context"
	"dataset/controller"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"strings"
	"testing"
)

const PlainTextEdit2AeneasScript = `is_new: yes
dataset_name: PlainTextEditScript2_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 14__plain_text_edit2_aeneas.csv
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps: 
  aeneas: yes
`

func TestPlainTextEdit2AeneasScript(t *testing.T) {
	type try struct {
		bibleId   string
		textNtId  string
		audioNTId string
		language  string
		expected  int
	}
	var tests []try
	tests = append(tests, try{bibleId: `ENGWEB`, textNtId: `ENGWEBN_ET`, audioNTId: `ENGWEBN2DA`,
		language: `eng`, expected: 8219})
	ctx := context.Background()
	for _, tst := range tests {
		var req = strings.Replace(PlainTextEdit2AeneasScript, `{bibleId}`, tst.bibleId, 2)
		var control = controller.NewController(ctx, []byte(req))
		filename, status := control.Process()
		if status.IsErr {
			t.Error(status)
		}
		fmt.Println(filename)
		numLines := NumCVSFileLines(filename, t)
		if numLines != tst.expected {
			t.Error(`Expected `, tst.expected, `records, got`, numLines)
		}
		identTest(`PlainTextEditScript2_`+tst.bibleId, t, request.TextPlainEdit, ``,
			tst.textNtId, ``, tst.audioNTId, tst.language)
	}
}

func TestCompareTimestamps(t *testing.T) {
	aenConn := openDatabase(`PlainTextEditScript2_ENGWEB`, t)
	bbConn := openDatabase(`PlainTextEditScript_ENGWEB`, t)
	var totalDiffs []float64
	for _, bookId := range db.RequestedBooks(request.Testament{NTBooks: []string{"1JN"}}) {
		lastChapter, _ := db.BookChapterMap[bookId]
		for chap := 1; chap <= lastChapter; chap++ {
			var diffs []float64
			bbTimestamps, status := bbConn.SelectScriptTimestamps(bookId, chap)
			if status.IsErr {
				t.Fatal(status)
			}
			var bbMap = make(map[string]float64)
			for _, ts := range bbTimestamps {
				bbMap[ts.VerseStr] = ts.BeginTS
			}
			aenTimestamps, status := aenConn.SelectScriptTimestamps(bookId, chap)
			if status.IsErr {
				t.Fatal(status)
			}
			for _, ts := range aenTimestamps {
				bbTS, ok := bbMap[ts.VerseStr]
				if !ok {
					t.Error(ts.VerseStr)
				}
				diff := ts.BeginTS - bbTS
				fmt.Println("verse:", ts.VerseStr, "AEN:", ts.BeginTS, "BB:", bbTS, "diff:", diff)
				diffs = append(diffs, diff)
			}
			mean, stddev := stat.MeanStdDev(diffs, nil)
			fmt.Println(bookId, chap, mean, stddev)
			totalDiffs = append(totalDiffs, diffs...)
		}
	}
	mean, stddev := stat.MeanStdDev(totalDiffs, nil)
	fmt.Println("Total", mean, stddev)
}

func openDatabase(datasetName string, t *testing.T) db.DBAdapter {
	ctx := context.Background()
	user, _ := fetch.GetTestUser()
	conn, status := db.NewerDBAdapter(ctx, false, user.Username, datasetName)
	if status.IsErr {
		t.Fatal(status)
	}
	return conn
}

// ENGWEB BB timestamps
// select avg(script_end_ts-script_begin_ts) from scripts where script_end_ts != 0.0
// = 8.37511692230324
