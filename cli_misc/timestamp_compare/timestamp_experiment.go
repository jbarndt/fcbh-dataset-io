package main

import (
	"context"
	"dataset"
	"dataset/cli_misc"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"fmt"
	"gonum.org/v1/gonum/stat"
)

type TSCompare struct {
	mediaId   string
	testament request.Testament
	apiTS     fetch.APIDBPTimestamps
	cliMisc   cli_misc.TSBucket
}

func main() {
	ctx := context.Background()
	cliMisc := cli_misc.NewTSBucket(ctx)
	//var testament = request.Testament{NTBooks: []string{"3JN"}}
	var testament = request.Testament{NT: true}
	tsDataPath := "cli_misc/find_timestamps/TestFilesetList.json"
	tsData := cliMisc.GetTSData(tsDataPath)
	var totalDiffs []float64
	for _, data := range tsData {
		diffs := CompareBB2Sandeep(cliMisc, data.MediaId, testament)
		totalDiffs = append(totalDiffs, diffs...)
	}
	mean, stddev := stat.MeanStdDev(totalDiffs, nil)
	fmt.Println("Final", mean, stddev)
}

func CompareBB2Sandeep(cliMisc cli_misc.TSBucket, mediaId string, testament request.Testament) []float64 {
	tsCompare := TSCompare{}
	tsCompare.mediaId = mediaId
	tsCompare.testament = testament
	tsCompare.apiTS = fetch.NewAPIDBPTimestamps(db.DBAdapter{}, mediaId)
	tsCompare.cliMisc = cliMisc
	return CompareFileset(tsCompare)
}

func CompareFileset(compare TSCompare) []float64 {
	var totalDiffs []float64
	for _, bookId := range db.RequestedBooks(compare.testament) {
		lastChapter, _ := db.BookChapterMap[bookId]
		for chap := 1; chap <= lastChapter; chap++ {
			timestamps1 := compare.cliMisc.GetTimestamps(cli_misc.VerseAeneas, compare.mediaId, bookId, chap)
			apiTS, status := compare.apiTS.Timestamps(bookId, chap)
			timestamps2 := ConvertAPI2DBTimestamp(apiTS)
			if status.IsErr {
				panic(status)
			}
			diffs := CompareChapter(compare.mediaId, bookId, chap, timestamps1, timestamps2)
			mean, stddev := stat.MeanStdDev(diffs, nil)
			fmt.Println(compare.mediaId, bookId, chap, mean, stddev)
			totalDiffs = append(totalDiffs, diffs...)
		}
	}
	mean, stddev := stat.MeanStdDev(totalDiffs, nil)
	fmt.Println("Total", compare.mediaId, mean, stddev)
	return totalDiffs
}

func CompareChapter(mediaId string, bookId string, chapter int, times1 []db.Timestamp, times2 []db.Timestamp) []float64 {
	ts1Map := make(map[int]db.Timestamp)
	for _, ts := range times1 {
		ts1Map[dataset.SafeVerseNum(ts.VerseStr)] = ts
	}
	var diffs []float64
	for _, ts2 := range times2 {
		ts1, ok := ts1Map[dataset.SafeVerseNum(ts2.VerseStr)]
		if !ok && ts2.VerseStr != `0` {
			fmt.Println("Not found ", mediaId, bookId, chapter, ts2.VerseStr)
			fmt.Println("times1", times1)
			fmt.Println("times2", times2)
		} else {
			diff := ts2.BeginTS - ts1.BeginTS
			//fmt.Println("verse:", ts2.VerseStr, "Base:", ts1.BeginTS, "Compare:", ts2.BeginTS, "diff:", diff)
			diffs = append(diffs, diff)
		}
	}
	return diffs
}

func ConvertAPI2DBTimestamp(ts []fetch.Timestamp) []db.Timestamp {
	var results []db.Timestamp
	for _, t := range ts {
		var dbTS db.Timestamp
		dbTS.VerseStr = t.VerseStart
		dbTS.BeginTS = t.Timestamp
		results = append(results, dbTS)
	}
	return results
}
