package fetch

import (
	"dataset"
	log "dataset/logger"
	"encoding/json"
	"fmt"
)

/*
This is not being used, because it does not recognize the limit parameter.
So, it would take 102 queries to get the entire list.
*/

type DownloadListType struct {
	Type         string `json:"type"`
	LanguageName string `json:"language"`
	Licensor     string `json:"licensor"`
	FilesetId    string `json:"fileset_id"`
}

type DownloadListResp struct {
	Data []DownloadListType `json:"data"`
	Meta any                `json:"meta"`
}

func (d *APIDBPClient) DownloadList() (map[string]DownloadListType, dataset.Status) {
	var result = make(map[string]DownloadListType)
	var status dataset.Status
	var get = `https://4.dbt.io/api/download/list?v=4`
	body, status := httpGet(d.ctx, get, false, d.bibleId)
	if status.IsErr {
		return result, status
	}
	fmt.Println(string(body))
	var response DownloadListResp
	err := json.Unmarshal(body, &response)
	if err != nil {
		status := log.Error(d.ctx, 500, err, "Error decoding DBP API /bibles JSON")
		return result, status
	}
	//fmt.Println(response)
	for _, item := range response.Data {
		result[item.FilesetId] = item
	}
	return result, status
}
