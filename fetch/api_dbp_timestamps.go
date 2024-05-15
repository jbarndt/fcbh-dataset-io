package fetch

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"dataset/request"
	"encoding/json"
	"strconv"
)

type APIDBPTimestamps struct {
	ctx     context.Context
	conn    db.DBAdapter
	audioId string
}

func NewAPIDBPTimestamps(conn db.DBAdapter, audioId string) APIDBPTimestamps {
	var a APIDBPTimestamps
	a.ctx = conn.Ctx
	a.conn = conn
	a.audioId = audioId
	return a
}

func (a *APIDBPTimestamps) LoadTimestamps(testament request.Testament) (bool, dataset.Status) {
	var result bool
	var status dataset.Status
	var audioIdMap map[string]bool
	audioIdMap, status = a.HavingTimestamps()
	if status.IsErr {
		return false, status
	}
	_, ok := audioIdMap[a.audioId]
	if !ok {
		status = log.ErrorNoErr(a.ctx, 400, `There are no timestamps available`)
		return false, status
	}
	scripts, status := a.conn.SelectScriptIds()
	if status.IsErr {
		return false, status
	}
	var scriptMap = make(map[string]int)
	for _, scp := range scripts {
		key := scp.BookId + ` ` + strconv.Itoa(scp.ChapterNum) + `:` + scp.VerseStr
		//_, _ = scriptMap[key]
		//if ok {
		//	log.Warn(a.ctx, `Duplicate book, chapter, verse`, key)
		//}
		scriptMap[key] = scp.ScriptId
	}
	lastBookId := ``
	lastChapter := -1
	for _, scp := range scripts {
		if scp.BookId != lastBookId || scp.ChapterNum != lastChapter {
			lastBookId = scp.BookId
			lastChapter = scp.ChapterNum
			if testament.HasNT(scp.BookId) || testament.HasOT(scp.BookId) {
				//fmt.Println("Getting Timestamps", scp.BookId, scp.ChapterNum)
				timestamp, status := a.Timestamps(scp.BookId, scp.ChapterNum)
				if status.IsErr {
					return false, status
				}
				var dbTimestamps []db.Timestamp
				var priorTS db.Timestamp
				for _, ts := range timestamp {
					var dbTS db.Timestamp
					dbTS.Id, ok = scriptMap[ts.Key()]
					if !ok {
						log.Warn(a.ctx, `Missing book, chapter, verse`, ts.Key())
					}
					dbTS.BeginTS = ts.Timestamp
					priorTS.EndTS = ts.Timestamp
					if priorTS.Id != 0 {
						dbTimestamps = append(dbTimestamps, priorTS)
					}
					priorTS = dbTS
				}
				dbTimestamps = append(dbTimestamps, priorTS)
				a.conn.UpdateScriptTimestamps(dbTimestamps)
				result = true
			}
		}
	}
	return result, status
}

func (a *APIDBPTimestamps) HavingTimestamps() (map[string]bool, dataset.Status) {
	var result = make(map[string]bool)
	var status dataset.Status
	var get = `https://4.dbt.io/api/timestamps?v=4`
	body, status := httpGet(a.ctx, get, false, `timestamps`)
	if status.IsErr {
		return result, status
	}
	var response []map[string]string
	err := json.Unmarshal(body, &response)
	if err != nil {
		status := log.Error(a.ctx, 500, err, "Error decoding DBP API /timestamp JSON")
		return result, status
	}
	for _, item := range response {
		for _, filesetId := range item {
			result[filesetId] = true
		}
	}
	return result, status
}

type Timestamp struct {
	BookId        string  `json:"book"`
	Chapter       string  `json:"chapter"`
	VerseStart    string  `json:"verse_start"`
	VerseStartAlt string  `json:"verse_start_alt"`
	Timestamp     float64 `json:"timestamp"`
}

func (t *Timestamp) Key() string {
	return t.BookId + ` ` + t.Chapter + `:` + t.VerseStart
}

type TimestampsResp struct {
	Data []Timestamp `json:"data"`
}

func (a *APIDBPTimestamps) Timestamps(bookId string, chapter int) ([]Timestamp, dataset.Status) {
	var result []Timestamp
	var status dataset.Status
	chapterStr := strconv.Itoa(chapter)
	var get = `https://4.dbt.io/api/timestamps/` + a.audioId + `/` + bookId + `/` + chapterStr + `?v=4`
	body, status := httpGet(a.ctx, get, false, `timestamps`)
	if status.IsErr {
		return result, status
	}
	//fmt.Println("BODY:", string(body))
	var response TimestampsResp
	err := json.Unmarshal(body, &response)
	if err != nil {
		status := log.Error(a.ctx, 500, err, "Error decoding DBP API /timestamp JSON")
		return result, status
	}
	result = response.Data
	return result, status
}
