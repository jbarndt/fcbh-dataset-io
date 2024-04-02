package fetch

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	HOST = "https://4.dbt.io/api/"
)

type DBPAPIClient struct {
	ctx     context.Context
	bibleId string
}

func NewDBPAPIClient(ctx context.Context, bibleId string) DBPAPIClient {
	var d DBPAPIClient
	d.ctx = ctx
	d.bibleId = bibleId
	return d
}

type FilesetType struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	Size      string `json:"size"`
	Codec     string `json:"codec"`
	Container string `json:"container"`
	Bitrate   string `json:"bitrate"`
}
type DbpProdType struct {
	Filesets []FilesetType `json:"dbp-prod"`
}
type AlphabetType struct {
	Alphabet string `json:"script"`
}
type BibleInfoType struct {
	BibleId       string `json:"abbr"`
	LanguageISO   string `json:"iso"`
	VersionCode   string
	LanguageId    int          `json:"language_id"`
	RolvId        int          `json:"language_rolv_code"`
	LanguageName  string       `json:"language"`
	VersionName   string       `json:"name"`
	Alphabet      AlphabetType `json:"alphabet"` // alphabet.script
	DbpProd       DbpProdType  `json:"filesets"`
	AudioFilesets []FilesetType
	TextFilesets  []FilesetType
}
type BibleInfoRespType struct {
	Data BibleInfoType `json:"data"`
}

func (d *DBPAPIClient) BibleInfo() (BibleInfoType, dataset.Status) {
	var result BibleInfoType
	var status dataset.Status
	result.VersionCode = d.bibleId[3:]
	var get = `https://4.dbt.io/api/bibles/` + d.bibleId + `?v=4`
	var response BibleInfoRespType
	body, status := d.httpGet(get, d.bibleId)
	if !status.IsErr {
		//if body != nil && len(body) > 0 {
		err := json.Unmarshal(body, &response)
		if err != nil {
			status := log.Error(d.ctx, 500, err, "Error decoding DBP API /bibles JSON")
			return result, status
		}
		result = response.Data
	}
	return result, status
}

func CreateIdent(info BibleInfoType) db.Ident {
	var id db.Ident
	id.BibleId = info.BibleId
	id.AudioFilesetId = ConcatFilesetId(info.AudioFilesets)
	id.TextFilesetId = ConcatFilesetId(info.TextFilesets)
	id.LanguageISO = info.LanguageISO
	id.VersionCode = info.VersionCode
	id.LanguageId = info.LanguageId
	id.RolvId = info.RolvId
	id.Alphabet = info.Alphabet.Alphabet
	id.LanguageName = info.LanguageName
	id.VersionName = info.VersionName
	return id
}

func (d *DBPAPIClient) FindFilesets(info *BibleInfoType, audio dataset.AudioSourceType,
	text dataset.TextSourceType, testament dataset.TestamentType) bool {
	var reqSize = string(testament)
	var okAudio = true
	var okText = true
	switch audio {
	case dataset.MP3:
		okAudio = d.searchAudio(info, `audio`, reqSize, string(dataset.MP3))
		if !okAudio {
			okAudio = d.searchAudio(info, `audio_drama`, reqSize, string(dataset.MP3))
		}
	}
	switch text {
	case dataset.SCRIPT:
		// audio_script
	case dataset.DBPTEXT:
		okText = d.searchText(info, `text_plain`, reqSize)
	case dataset.TEXTEDIT:
		okText = d.searchText(info, `text_plain`, reqSize)
	case dataset.USXEDIT:
		okText = d.searchText(info, `text_usx`, reqSize)
	}
	return okAudio && okText
}

func (d *DBPAPIClient) searchText(info *BibleInfoType, reqType string, reqSize string) bool {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == reqType && (reqSize == `C` || rec.Size == reqSize) {
			info.TextFilesets = append(info.TextFilesets, rec)
		}
	}
	return len(info.TextFilesets) > 0
}

func (d *DBPAPIClient) searchAudio(info *BibleInfoType, reqType string, reqSize string, reqCodec string) bool {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == reqType && (reqSize == `C` || rec.Size == reqSize) {
			if strings.ToUpper(rec.Codec) == reqCodec {
				info.AudioFilesets = append(info.AudioFilesets, rec)
			}
		}
	}
	return len(info.AudioFilesets) > 0
}

func ConcatFilesetId(filesets []FilesetType) string {
	var result []string
	for _, rec := range filesets {
		result = append(result, rec.Id)
	}
	return strings.Join(result, `,`)
}

func (d *DBPAPIClient) Download(info BibleInfoType) dataset.Status {
	var status dataset.Status
	var directory = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), info.BibleId)
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return log.Error(d.ctx, 500, err, `Could not create directory to store downloaded files.`)
		}
	}
	var download []FilesetType
	download = append(download, info.TextFilesets...)
	download = append(download, info.AudioFilesets...)
	for _, rec := range download {
		if rec.Type == `text_plain` {
			status = d.downloadPlainText(directory, rec.Id)
			if status.IsErr {
				return status
			}
		} else {
			locations, status := d.downloadLocation(rec.Id)
			if status.IsErr {
				return status
			}
			locations, status = d.sortFileLocations(locations)
			if status.IsErr {
				return status
			}
			directory = filepath.Join(directory, rec.Id)
			status = d.downloadFiles(directory, locations)
			if status.IsErr {
				return status
			}
		}
	}
	return status
}

func (d *DBPAPIClient) downloadPlainText(directory string, filesetId string) dataset.Status {
	var content []byte
	var status dataset.Status
	filename := filesetId + ".json"
	filePath := filepath.Join(directory, filename)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		var get = HOST + "download/" + filesetId + "?v=4&limit=100000"
		fmt.Println("Downloading to", filePath)
		content, status = d.httpGet(get, filesetId)
		if !status.IsErr {
			d.saveFile(filePath, content)
		}
	}
	return status
}

type LocationRec struct {
	BookId   string `json:"book_id"`
	BookName string `json:"book_name"`
	Chapter  int    `json:"chapter_start"`
	Verse    int    `json:"verse_start"`
	URL      string `json:"path"`
	FileSize int    `json:"filesize_in_bytes"`
	Filename string
}
type LocationDownloadRec struct {
	Data []LocationRec `json:"data"`
	Meta any           `json:"meta"`
}

func (d *DBPAPIClient) downloadLocation(filesetId string) ([]LocationRec, dataset.Status) {
	var result []LocationRec
	var status dataset.Status
	var get string
	if strings.Contains(filesetId, `usx`) {
		get = HOST + "bibles/filesets/" + filesetId + "/ALL/1?v=4&limit=100000"
	} else {
		get = HOST + "download/" + filesetId + "?v=4"
	}
	var content []byte
	content, status = d.httpGet(get, filesetId)
	//if len(content) == 0 {
	if !status.IsErr {
		var response LocationDownloadRec
		err := json.Unmarshal(content, &response)
		if err != nil {
			status = log.Error(d.ctx, 500, err, "Error parsing json for", filesetId)
		} else {
			result = response.Data
		}
	}
	return result, status
}

func (d *DBPAPIClient) sortFileLocations(locations []LocationRec) ([]LocationRec, dataset.Status) {
	var status dataset.Status
	for i, loc := range locations {
		get, err := url.Parse(loc.URL)
		if err != nil {
			status = log.Error(d.ctx, 500, err, "Could not parse URL", loc.URL)
			if status.IsErr {
				return locations, status
			}
		}
		locations[i].Filename = filepath.Base(get.Path)
	}
	sort.Slice(locations, func(i int, j int) bool {
		return locations[i].Filename < locations[j].Filename
	})
	return locations, status
}

func (d *DBPAPIClient) downloadFiles(directory string, locations []LocationRec) dataset.Status {
	var status dataset.Status
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return log.Error(d.ctx, 500, err, "Could not create directory to store downloaded files.")
		}
	}
	for _, loc := range locations {
		size := loc.FileSize
		filePath := filepath.Join(directory, loc.Filename)
		file, err := os.Stat(filePath)
		if os.IsNotExist(err) || file.Size() != int64(size) {
			fmt.Println("Downloading", loc.Filename)
			var content []byte
			content, status = d.httpGet(loc.URL, loc.Filename)
			if !status.IsErr {
				if len(content) != size {
					log.Warn(d.ctx, "Warning for", loc.Filename, "has an expected size of", size, "but, actual size is", len(content))
				}
				d.saveFile(filePath, content)
			}
		}
	}
	return status
}

func (d *DBPAPIClient) httpGet(url string, desc string) ([]byte, dataset.Status) {
	var body []byte
	var status dataset.Status
	url += `&limit=100000&key=` + os.Getenv(`FCBH_DBP_KEY`)
	resp, err := http.Get(url)
	if err != nil {
		status = log.Error(d.ctx, resp.StatusCode, err, "Error in DBP API request for:", desc)
		return body, status
	}
	defer resp.Body.Close()
	if resp.Status[0] == '2' {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			status = log.Error(d.ctx, resp.StatusCode, err, "Error reading DBP API response for:", desc)
			return body, status
		}
	}
	return body, status
}

func (d *DBPAPIClient) saveFile(filePath string, content []byte) dataset.Status {
	var status dataset.Status
	fp, err := os.Create(filePath)
	if err != nil {
		return log.Error(d.ctx, 500, err, "Error Creating file during download.")
	}
	_, err = fp.Write(content)
	if err != nil {
		return log.Error(d.ctx, 500, err, "Error writing to file during download.")
	}
	err = fp.Close()
	if err != nil {
		return log.Error(d.ctx, 500, err, "Error closing file during download.")
	}
	return status
}
