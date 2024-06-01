package testing

import (
	"context"
	"dataset/controller"
	"dataset/db"
	"dataset/input"
	"dataset/request"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"strings"
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

func TestPlainTextAeneasTimestampsScript(t *testing.T) {
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 8219, TextNtId: "ENGWEBN_ET",
		AudioNTId: "ENGWEBN2DA", Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(PlainTextEditBBTimestampsScript, tests, t)
}

func TestCompareTimestamps(t *testing.T) {
	ctx := context.Background()
	aws := input.NewTSBucket(ctx)
	tsData := aws.GetTSData()
	count := 0
	for _, tst := range tsData {
		if count > 0 {
			break
		}
		count++
		var req = strings.Replace(PlainTextEditBBTimestampsScript, `{bibleId}`, tst.MediaId[:6], 2)
		username, datasetName := extractKeyFields(ctx, req)
		if !db.DatabaseExists(username, datasetName) {
			var control = controller.NewController(ctx, []byte(req))
			filename, status := control.Process()
			if status.IsErr {
				t.Fatal(status)
			}
			fmt.Println(filename)
			//numLines := NumFileLines(filename, t)
		}
		conn, status := db.NewerDBAdapter(ctx, false, username, datasetName)
		if status.IsErr {
			t.Fatal(status)
		}
		var totalDiffs []float64
		for _, bookId := range db.RequestedBooks(request.Testament{NTBooks: []string{"1JN"}}) {
			lastChapter, _ := db.BookChapterMap[bookId]
			for chap := 1; chap <= lastChapter; chap++ {
				var diffs []float64
				timestamps := aws.GetTimestamps(input.VerseAeneas, tst.MediaId, bookId, chap)
				var baseMap = make(map[string]float64)
				for _, ts := range timestamps {
					baseMap[ts.VerseStr] = ts.BeginTS
				}
				datasetTS, status := conn.SelectScriptTimestamps(bookId, chap)
				if status.IsErr {
					t.Fatal(status)
				}
				for _, ts := range datasetTS {
					baseTS, ok := baseMap[ts.VerseStr]
					if !ok {
						t.Error(ts.VerseStr + ` is not found in baseTS`)
					}
					diff := ts.BeginTS - baseTS
					fmt.Println("verse:", ts.VerseStr, "Base:", baseTS, "Dataset:", ts.BeginTS, "diff:", diff)
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
}

func extractKeyFields(ctx context.Context, requestYaml string) (string, string) {
	decoder := request.NewRequestDecoder(ctx)
	req, status := decoder.Decode([]byte(requestYaml))
	if status.IsErr {
		panic(status)
	}
	return req.Username, req.DatasetName
}
