package main

import (
	"context"
	"dataset/db"
	"dataset/decode_yaml/request"
	"dataset/fetch"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type TS struct {
	Book        string  `json:"book"`
	Chap        int     `json:"chap"`
	Verse       string  `json:"verse"`
	BBBeginTS   float64 `json:"bb_begin_ts"`
	BBAbsent    bool    `json:"bb_absent"`
	WahaBeginTS float64 `json:"waha_begin_ts"`
	WahaEndTS   float64 `json:"waha_end_ts"`
	BeginTSDiff float64 `json:"begin_diff"`
}

func main() {
	directory := filepath.Join(os.Getenv("FCBH_DATASET_DB"), "Muller_Timestamps")
	//LoadBBTimestamps(directory, "NPIDPI", "NPIDPIN1DA")
	//LoadBBTimestamps(directory, "URDIRV", "URDIRVN1DA")
	//ConpareWahaTS2BBTS(directory, "npi", "NPIDPIN1DA")
	//ConpareWahaTS2BBTS(directory, "urd", "URDIRVN1DA")
	ComputeAverageAnDSort(directory, "npi")
	//ComputeAverageAnDSort(directory, "urd")
	//CompareBB2S3TS(directory, "npi", "NPIDPIN1DA")
}

// LoadBBTimestamps is is a utility for getting timestamp data from API into a json file.
func LoadBBTimestamps(directory string, bibleId string, filesetId string) {
	var results []fetch.Timestamp
	ctx := context.Background()
	database := bibleId + `_TS`
	conn, status := db.NewerDBAdapter(ctx, true, `GaryNGriswold`, database)
	if status != nil {
		panic(status)
	}
	api := fetch.NewAPIDBPTimestamps(conn, filesetId)
	books := db.RequestedBooks(request.Testament{OT: false, NT: true})
	for _, book := range books {
		maxChap, _ := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			fmt.Println(book, chap)
			tsList, status := api.Timestamps(book, chap)
			if status != nil {
				panic(status)
			}
			//fmt.Println(book, chap, tsList)
			results = append(results, tsList...)
		}
	}
	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(bytes))
	filename := filepath.Join(directory, "BB_"+strings.ToLower(bibleId[:3])+".json")
	os.WriteFile(filename, bytes, 0644)
}

func ConpareWahaTS2BBTS(directory string, isoCode string, mediaId string) {
	var results []TS
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, ":memory:")
	api := fetch.NewAPIDBPTimestamps(conn, mediaId)
	books := db.RequestedBooks(request.Testament{OT: false, NT: true}) //, NTBooks: []string{`MRK`}})
	//request.Testament{NTBooks: []string{}}
	for _, book := range books {
		maxChap, _ := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			fmt.Println(book, chap)
			bbTSMap := readBBTS(api, book, chap)
			wahaChap := readWaha(directory, isoCode, book, chap)
			//fmt.Println(book, chap, wahaChap)
			for _, wahaVerse := range wahaChap.Verses {
				var ts TS
				var ok bool
				ts.Book = book
				ts.Chap = chap
				parts := strings.Split(wahaVerse.VerseId, ".")
				ts.Verse = parts[2]
				ts.WahaBeginTS = wahaVerse.Timings[0]
				ts.WahaEndTS = wahaVerse.Timings[1]
				ts.BBBeginTS, ok = bbTSMap[ts.Verse]
				if !ok {
					ts.BBAbsent = true
					fmt.Println("Missing bb_begin_ts", wahaVerse)
				} else {
					ts.BeginTSDiff = ts.BBBeginTS - ts.WahaBeginTS
				}
				results = append(results, ts)
			}
		}
	}
	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(bytes))
	filename := filepath.Join(directory, isoCode+"_compare.json")
	os.WriteFile(filename, bytes, 0644)
}

type WahaChap struct {
	TranslationId string      `json:"translationId"`
	BookName      string      `json:"bookName"`
	ChapterId     string      `json:"chapterId"`
	Reference     string      `json:"reference"`
	Verses        []WahaVerse `json:"verses"`
}

type WahaVerse struct {
	VerseId    string    `json:"verseId"`
	Text       string    `json:"text"`
	Timings    []float64 `json:"timings"`
	TimingsStr []string  `json:"timings_str"`
	URoman     string    `json:"uroman"`
}

func readWaha(directory string, isoCode string, book string, chap int) *WahaChap {
	var result WahaChap
	chapStr := fmt.Sprintf("%03s", strconv.Itoa(chap))
	filename := book + "_" + chapStr + ".json"
	chapFile := filepath.Join(directory, isoCode, book, filename)
	//fmt.Println(chapFile)
	bytes, err := os.ReadFile(chapFile)
	if err != nil {
		fmt.Println("File not found", chapFile)
	} else {
		json.Unmarshal(bytes, &result)
	}
	return &result
}

func readBBTS(api fetch.APIDBPTimestamps, book string, chap int) map[string]float64 {
	var result = make(map[string]float64)
	timestamps, status := api.Timestamps(book, chap)
	if status != nil {
		panic(status)
	}
	//fmt.Println(timestamps)
	for _, ts := range timestamps {
		result[ts.VerseStart] = ts.Timestamp
	}
	return result
}

func readS3TS(directory string, isoCode string, book string, chap int) map[string]float64 {
	var result = make(map[string]float64)
	bookSeq, ok := db.BookSeqMap[book]
	if !ok {
		panic("Missing book " + book)
	}
	bookSeq--
	chapStr := fmt.Sprintf("%02s", strconv.Itoa(chap))
	filename := "C01-" + strconv.Itoa(bookSeq) + "-" + book + "-" + chapStr + "-timing.txt"
	//var filename = "C01-66-REV-17-timing.txt"
	filePath := filepath.Join(directory, isoCode+"_S3", filename)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) == 3 {
			verse := parts[2]
			beginTS, err1 := strconv.ParseFloat(parts[0], 64)
			if err1 != nil {
				panic(err1)
			}
			result[verse] = beginTS
		}
	}
	return result
}

func ComputeAverageAnDSort(directory string, isoCode string) {
	filename := filepath.Join(directory, isoCode+"_compare.json")
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var timestamps []TS
	err = json.Unmarshal(bytes, &timestamps)
	if err != nil {
		panic(err)
	}
	var diffs []float64
	for _, ts := range timestamps {
		//if ts.Verse != "1" && ts.BBAbsent == false {
		if ts.BBAbsent == false {
			diffs = append(diffs, ts.BeginTSDiff)
		}
	}
	fmt.Printf("data sample size: %v\n", len(diffs))
	mean := stat.Mean(diffs, nil)
	variance := stat.Variance(diffs, nil)
	stddev := math.Sqrt(variance)
	sort.Float64s(diffs)
	median := stat.Quantile(0.5, stat.Empirical, diffs, nil)
	fmt.Printf("mean=     %v\n", mean)
	fmt.Printf("median=   %v\n", median)
	fmt.Printf("variance= %v\n", variance)
	fmt.Printf("std-dev=  %v\n", stddev)
	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].BeginTSDiff > timestamps[j].BeginTSDiff
	})
	var count = 0
	for _, ts := range timestamps {
		if ts.Verse != "1" && ts.BBAbsent == false {
			fmt.Printf("%d\t%g\t%+v\n", count, ts.BeginTSDiff, ts)
			count++
		}
		if count >= 16 {
			break
		}
	}
}

func CompareBB2S3TS(directory string, isoCode string, mediaId string) {
	var diffs []float64
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, ":memory:")
	api := fetch.NewAPIDBPTimestamps(conn, mediaId)
	for _, book := range db.RequestedBooks(request.Testament{NT: true}) {
		maxChap, _ := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			fmt.Println(book, chap)
			s3TSMap := readS3TS(directory, isoCode, book, chap)
			bbTSMap := readBBTS(api, book, chap)
			for verse, beginS3TS := range s3TSMap {
				beginBBTS, ok := bbTSMap[verse]
				if !ok {
					fmt.Println("Missing BBTS", book, chap, verse)
				} else {
					diff := beginBBTS - beginS3TS
					diffs = append(diffs, diff)
				}
			}
		}
	}
	mean := stat.Mean(diffs, nil)
	variance := stat.Variance(diffs, nil)
	stddev := math.Sqrt(variance)
	sort.Float64s(diffs)
	median := stat.Quantile(0.5, stat.Empirical, diffs, nil)
	fmt.Printf("mean=     %v  %.4f\n", mean, mean)
	fmt.Printf("median=   %v  %.4f\n", median, median)
	fmt.Printf("variance= %v  %.4f\n", variance, variance)
	fmt.Printf("std-dev=  %v  %.4f\n", stddev, stddev)
}
