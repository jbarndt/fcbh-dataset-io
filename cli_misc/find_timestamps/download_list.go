package main

import (
	"context"
	"dataset"
	"dataset/fetch"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

/*
Write method to do /download/list

Write a program that reads Sandeeps bucket in sequence,

1. Timestamps by script
2. Timestamps by verse
3. Script

If useful /timestamps find those with timestamps

Produce a file of those that have all resources, with path for each thing.
Possibly json or csv output.
*/

type TSData struct {
	MediaId string `json:"media_id"`
}

type DownloadList struct {
	TType     string `json:"type"`
	Language  string `json:"language"`
	FilesetId string `json:"fileset_id"`
}
type Pagination struct {
	Total       int `json:"total"`
	Count       int `json:"count"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}
type Pagination1 struct {
	Page Pagination `json:"pagination"`
}
type DownloadResp struct {
	Data []DownloadList `json:"data"`
	Meta Pagination1    `json:"meta"`
}

func main() {
	downloadList := downloadFilestList()
	fmt.Println(downloadList)
}

func downloadFilestList() []TSData {
	var results []TSData
	ctx := context.Background()
	var page = 0
	for {
		page++
		url := fetch.HOST + `download/list?page=` + strconv.Itoa(page) + `&v=4`
		content, status := fetch.HttpGet(ctx, url, `fetch allowed filesets`)
		catchStatus(status)
		//fmt.Println(string(content))
		var list DownloadResp
		err := json.Unmarshal(content, &list)
		catchErr(err)
		fmt.Println(list.Meta)
		for _, item := range list.Data {
			if item.TType == `audio` || item.TType == `audio_drama` {
				var rec TSData
				rec.MediaId = item.FilesetId
				results = append(results, rec)
			}
		}
		if list.Meta.Page.CurrentPage == list.Meta.Page.TotalPages {
			break
		}
	}
	//var jsonContent []byte
	jsonContent, err := json.MarshalIndent(results, "", "    ")
	catchErr(err)
	err = os.WriteFile(`cli_misc/find_timestamps/FilesetList.json`, jsonContent, 0644)
	catchErr(err)
	return results
}

func catchErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func catchStatus(status dataset.Status) {
	if status.IsErr {
		fmt.Println(status)
		os.Exit(1)
	}
}
