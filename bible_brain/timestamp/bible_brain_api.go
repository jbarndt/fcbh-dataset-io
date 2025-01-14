package timestamp

import (
	"context"
	"dataset"
	"dataset/fetch"
	log "dataset/logger"
	"encoding/json"
	"os"
	"strings"
)

type BiblesResponse struct {
	Data []Bible `json:"data"`
	Meta Meta    `json:"meta"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	CurrentPage int    `json:"current_page"`
	FromPage    int    `json:"from"`
	LastPage    int    `json:"last_page"`
	NextPageURL string `json:"next_page_url"`
	PerPage     int    `json:"per_page"`
	PrevPageURL string `json:"prev_page_url"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

type Bible struct {
	Abbr       string `json:"abbr"`
	Name       string `json:"name"`
	Vname      string `json:"vname"`
	Autonym    string `json:"autonym"`
	LanguageId int    `json:"language_id"`
	Rolv_code  string `json:"rolv_code"`
	Language   string `json:"language"`
	ISOCode    string `json:"iso"`
	Date       string `json:"date"`
	Filesets   Bucket `json:"filesets"`
}

type Bucket struct {
	DBPProd []Fileset `json:"dbp-prod"`
	DBPVid  []Fileset `json:"dbp-vid"`
}

type Fileset struct {
	FilesetId string `json:"id"`
	Type      string `json:"type"`
	Size      string `json:"size"`
	StockNo   string `json:"stock_no"`
	Bitrate   string `json:"bitrate"`
	Codec     string `json:"codec"`
	Container string `json:"container"`
	Volume    string `json:"volume"`
}

func FetchBibles() dataset.Status {
	var result []Bible
	var status dataset.Status
	ctx := context.Background()
	url := "https://4.dbt.io/api/bibles?v=4"
	for {
		if url == "" {
			break
		}
		content, status := fetch.HttpGet(ctx, url, "")
		if status.IsErr {
			return status
		}
		var response BiblesResponse
		err := json.Unmarshal(content, &response)
		if err != nil {
			return log.Error(ctx, 500, err, "Error unmarshalling bibles API query")
		}
		result = append(result, response.Data...)
		url = response.Meta.Pagination.NextPageURL
		url = strings.Replace(url, "?page", "?v=4&page", -1)
	}
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return log.Error(ctx, 500, err, "Error marshalling bibles API query")
	}
	filename := "bible_brain_api.json"
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return log.Error(ctx, 500, err, "Error writing bibles API query")
	}
	return status
}
