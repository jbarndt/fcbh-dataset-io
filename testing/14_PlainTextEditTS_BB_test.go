package testing

import (
	"context"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"fmt"
	"gonum.org/v1/gonum/stat"
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
		AudioNTId: "ENGWEBN2DA", Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(PlainTextEditTSBBScript, tests, t)
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
