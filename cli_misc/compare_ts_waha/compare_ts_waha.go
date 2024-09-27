package main

import (
	"context"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	//LoadTimestamps()

	ConpareWahaTS2BBTS("npi", "NPIDPIN1DA")
}

// LoadTimestamps is is a utility for getting timestamp data from API into a json file.
func LoadTimestamps() {
	var results []fetch.Timestamp
	ctx := context.Background()
	bibleId := `NPIDPI`
	//bibleId := `URDIRV`
	database := bibleId + `_TS`
	conn, status := db.NewerDBAdapter(ctx, true, `GaryNGriswold`, database)
	if status.IsErr {
		panic(status)
	}
	filesetId := `NPIDPIN1DA`
	//filesetId := `URDIRVN1DA`
	api := fetch.NewAPIDBPTimestamps(conn, filesetId)
	books := db.RequestedBooks(request.Testament{OT: false, NT: true})
	for _, book := range books {
		maxChap, _ := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			tsList, status := api.Timestamps(book, chap)
			if status.IsErr {
				panic(status.Message)
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
	filename := filepath.Join(os.Getenv("FCBH_DATASET_DB"), "Muller_Timestamps", "npi.json")
	//filename := filepath.Join(os.Getenv("FCBH_DATASET_DB"), "Muller_Timestamps", "urd.json")
	os.WriteFile(filename, bytes, 0644)
}

func ConpareWahaTS2BBTS(isoCode string, mediaId string) {
	directory := filepath.Join(os.Getenv("FCBH_DATASET_DB"), "Muller_Timestamps")
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
	var results []TS
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, ":memory:")
	api := fetch.NewAPIDBPTimestamps(conn, mediaId)
	books := db.RequestedBooks(request.Testament{OT: false, NT: true}) //, NTBooks: []string{`MRK`}})
	//request.Testament{NTBooks: []string{}}
	for _, book := range books {
		fmt.Println(book)
		maxChap, _ := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			fmt.Println(chap)
			//maxChap = 2
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
	chapStr := fmt.Sprintf("%03s", strconv.Itoa(chap))
	filename := book + "_" + chapStr + ".json"
	chapFile := filepath.Join(directory, isoCode, book, filename)
	//fmt.Println(chapFile)
	bytes, err := os.ReadFile(chapFile)
	if err != nil {
		panic(err)
	}
	var result WahaChap
	json.Unmarshal(bytes, &result)
	return &result
}

func readBBTS(api fetch.APIDBPTimestamps, book string, chap int) map[string]float64 {
	var result = make(map[string]float64)
	timestamps, status := api.Timestamps(book, chap)
	if status.IsErr {
		panic(status)
	}
	//fmt.Println(timestamps)
	for _, ts := range timestamps {
		result[ts.VerseStart] = ts.Timestamp
	}
	return result
}
