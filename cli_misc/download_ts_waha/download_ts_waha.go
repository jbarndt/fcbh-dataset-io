package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/db"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/decode_yaml/request"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/fetch"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/input"
	log "github.com/faithcomesbyhearing/fcbh-dataset-io/logger"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/read"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	bibleId := "ABIWBT"
	outDir := filepath.Join(os.Getenv("HOME"), "miniconda3", "envs", "torch2", "data", bibleId)
	err := os.MkdirAll(outDir, 0777)
	if err != nil {
		panic(err)
	}
	testament := request.Testament{NTBooks: []string{"MRK"}}
	testament.BuildBookMaps()
	var info = downloadBible(bibleId, testament)
	renameAudioFiles(outDir, bibleId, info)
	chopupTextFile2Txt(outDir, bibleId, info, testament)
}

func downloadBible(bibleId string, testament request.Testament) fetch.BibleInfoType {
	var info fetch.BibleInfoType
	var status *log.Status
	ctx := context.Background()
	client := fetch.NewAPIDBPClient(ctx, bibleId)
	info, status = client.BibleInfo()
	if status != nil {
		panic(status)
	}
	audioData := request.BibleBrainAudio{MP3_64: true}
	textData := request.BibleBrainText{TextPlain: true}
	client.FindFilesets(&info, audioData, textData, testament)
	download := fetch.NewAPIDownloadClient(ctx, bibleId, testament)
	status = download.Download(info)
	if status != nil {
		panic(status)
	}
	return info
}

func renameAudioFiles(outputDir string, bibleId string, info fetch.BibleInfoType) {
	ctx := context.Background()
	dir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, info.AudioNTFileset.Id)
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		var inp input.InputFile
		inp.Filename = file.Name()
		inp.Directory = dir
		inp.MediaType = request.Audio
		input.ParseFilenames(ctx, &inp)
		fmt.Println(inp.BookId, inp.Chapter, inp.FilePath())
		newFilename := inp.BookId + "." + strconv.Itoa(inp.Chapter) + ".mp3"
		newFilePath := filepath.Join(outputDir, newFilename)
		os.Rename(inp.FilePath(), newFilePath)
	}
}

type PlainText struct {
	Verse string `json:"verse"`
	Text  string `json:"text"`
}

func chopupTextFile2Txt(outputDir string, bibleId string, info fetch.BibleInfoType,
	testament request.Testament) {
	ctx := context.Background()
	dir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	filename := filepath.Join(dir, info.AudioNTFileset.Id)
	fmt.Println(filename)
	var conn = db.NewDBAdapter(ctx, ":memory:")
	var reader = read.NewDBPTextReader(conn, testament)
	var files []input.InputFile
	var file input.InputFile
	file.Directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	file.Filename = info.TextNTPlainFileset.Id + ".json"
	file.MediaType = request.TextPlain
	files = append(files, file)
	status := reader.ProcessFiles(files)
	if status != nil {
		panic(status)
	}
	for _, book := range db.RequestedBooks(testament) {
		maxChap := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			verses, status2 := conn.SelectScriptsByChapter(book, chap)
			if status2 != nil {
				panic(status2)
			}
			var results []string
			for _, verse := range verses {
				//fmt.Printf("%+v\n", verse)
				results = append(results, verse.ScriptText)
			}
			bytes := []byte(strings.Join(results, "\n"))
			newPath := filepath.Join(outputDir, book+"."+strconv.Itoa(chap)+".txt")
			err := os.WriteFile(newPath, bytes, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}

// chopupTextFile2Json did not work. I am not sure what is the correct json format.
func chopupTextFile2Json(outputDir string, bibleId string, info fetch.BibleInfoType,
	testament request.Testament) {
	ctx := context.Background()
	dir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	filename := filepath.Join(dir, info.AudioNTFileset.Id)
	fmt.Println(filename)
	var conn = db.NewDBAdapter(ctx, ":memory:")
	var reader = read.NewDBPTextReader(conn, testament)
	var files []input.InputFile
	var file input.InputFile
	file.Directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
	file.Filename = info.TextNTPlainFileset.Id + ".json"
	file.MediaType = request.TextPlain
	files = append(files, file)
	status := reader.ProcessFiles(files)
	if status != nil {
		panic(status)
	}
	for _, book := range db.RequestedBooks(testament) {
		maxChap := db.BookChapterMap[book]
		for chap := 1; chap <= maxChap; chap++ {
			verses, status2 := conn.SelectScriptsByChapter(book, chap)
			if status2 != nil {
				panic(status2)
			}
			var results []PlainText
			for _, verse := range verses {
				//fmt.Printf("%+v\n", verse)
				var rec PlainText
				rec.Verse = verse.VerseStr
				rec.Text = verse.ScriptText
				results = append(results, rec)
			}
			bytes, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				panic(err)
			}
			newPath := filepath.Join(outputDir, book+"."+strconv.Itoa(chap)+".json")
			err = os.WriteFile(newPath, bytes, 0644)
		}
	}
}
