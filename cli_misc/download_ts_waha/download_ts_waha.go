package main

import (
	"context"
	"dataset"
	"dataset/fetch"
	"dataset/request"
)

func main() {
	bibleId := "ABIWBT"
	//directory := filepath.Join(os.Getenv("HOME"), "miniconda3", "envs", "waha_ts", "data")
	downloadBibleInfo(bibleId)
}

// func GetAllBibleId()

// func CheckLangTreeForSupport()

func downloadBibleInfo(bibleId string) {
	ctx := context.Background()
	client := fetch.NewAPIDBPClient(ctx, bibleId)
	var info fetch.BibleInfoType
	var status dataset.Status
	info, status = client.BibleInfo()
	if status.IsErr {
		panic(status)
	}
	audioData := request.BibleBrainAudio{MP3_64: true}
	textData := request.BibleBrainText{TextPlain: true}
	testament := request.Testament{NTBooks: []string{"MRK"}}
	testament.BuildBookMaps()
	client.FindFilesets(&info, audioData, textData, testament)
	download := fetch.NewAPIDownloadClient(ctx, bibleId, testament)
	status = download.Download(info)
	if status.IsErr {
		panic(status)
	}
}
