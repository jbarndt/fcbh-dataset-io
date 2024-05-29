package testing

import (
	"context"
	"dataset/cli_misc"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gonum.org/v1/gonum/stat"
	"os"
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
	type testCase struct {
		dataset string
		audioId string
	}
	var tests []testCase
	var test = testCase{dataset: `PlainTextEditTSScript_ENGWEB`, audioId: `ENGWEBN2DA`}
	tests = append(tests, test)
	ctx := context.Background()
	client := openAwsS3(ctx, t)
	content, err := os.ReadFile("../cli_misc/find_timestamps/TestFilesetList.json")
	if err != nil {
		t.Fatal(err)
	}
	var tsData []cli_misc.TSData
	err = json.Unmarshal(content, &tsData)
	if err != nil {
		t.Fatal(err)
	}
	var tsMap = make(map[string]cli_misc.TSData)
	for _, ts := range tsData {
		tsMap[ts.MediaId] = ts
	}
	for _, tst := range tests {
		conn := openDatabase(tst.dataset, t)
		var totalDiffs []float64
		for _, bookId := range db.RequestedBooks(request.Testament{NTBooks: []string{"1JN"}}) {
			lastChapter, _ := db.BookChapterMap[bookId]
			for chap := 1; chap <= lastChapter; chap++ {
				//var diffs []float64
				datasetTS, status := conn.SelectScriptTimestamps(bookId, chap)
				if status.IsErr {
					t.Fatal(status)
				}
				var datasetMap = make(map[string]float64)
				for _, ts := range datasetTS {
					datasetMap[ts.VerseStr] = ts.BeginTS
				}
				ts, ok := tsMap[tst.audioId]
				if !ok {
					t.Fatal(tst.audioId + ` is not found.`)
				}

				keys := listFiles(client, `dbp-aeneas-staging`, ts.ScriptTSPath, t)
				for _, obj := range keys {
					fmt.Println(obj)
				}

				//for _, ts := range aenTimestamps {
				//	bbTS, ok := bbMap[ts.VerseStr]
				//	if !ok {
				//		t.Error(ts.VerseStr)
				//	}
				//	diff := ts.BeginTS - bbTS
				//	fmt.Println("verse:", ts.VerseStr, "AEN:", ts.BeginTS, "BB:", bbTS, "diff:", diff)
				//	diffs = append(diffs, diff)
				//}
				//mean, stddev := stat.MeanStdDev(diffs, nil)
				//fmt.Println(bookId, chap, mean, stddev)
				//totalDiffs = append(totalDiffs, diffs...)
			}
		}
		mean, stddev := stat.MeanStdDev(totalDiffs, nil)
		fmt.Println("Total", mean, stddev)
	}
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

func openAwsS3(ctx context.Context, t *testing.T) *s3.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "us-west-2"
	})
	return client
}

func listFiles(client *s3.Client, bucket string, prefix string, t *testing.T) []string {
	ctx := context.Background()
	list, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
		//Delimiter: aws.String("/"),
	})
	if err != nil {
		t.Fatal(err)
	}
	var results []string
	for _, obj := range list.Contents {
		results = append(results, *obj.Key)
	}
	return results
}
