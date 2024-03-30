package fetch

import (
	"dataset"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	bibleId string
	//audioSource dataset_io.AudioSourceType
}

func NewDBPAPIClient(bibleId string) DBPAPIClient {
	var d DBPAPIClient
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

func (d *DBPAPIClient) BibleInfo() BibleInfoType {
	var result BibleInfoType
	var url = `https://4.dbt.io/api/bibles/` + d.bibleId + `?v=4`
	var response BibleInfoRespType
	body, status := d.httpGet(url, d.bibleId)
	if body != nil && len(body) > 0 {
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error decoding DBP API JSON:", status, err)
			return BibleInfoType{}
		}
		result = response.Data
	}
	result.VersionCode = d.bibleId[3:]
	return result
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

func (d *DBPAPIClient) Download(info BibleInfoType) {
	var directory = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), info.BibleId)
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	var download []FilesetType
	download = append(download, info.TextFilesets...)
	download = append(download, info.AudioFilesets...)
	for _, rec := range download {
		if rec.Type == `text_plain` {
			d.downloadPlainText(directory, rec.Id)
		} else {
			locations := d.downloadLocation(rec.Id)
			locations = d.sortFileLocations(locations)
			directory = filepath.Join(directory, rec.Id)
			d.downloadFiles(directory, locations)
		}
	}
}

func (d *DBPAPIClient) downloadPlainText(directory string, filesetId string) {
	filename := filesetId + ".json"
	filePath := filepath.Join(directory, filename)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		var url = HOST + "download/" + filesetId + "?v=4&limit=100000"
		fmt.Println("Downloading to", filePath)
		content, _ := d.httpGet(url, filesetId)
		d.saveFile(filePath, content)
	}
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

func (d *DBPAPIClient) downloadLocation(filesetId string) []LocationRec {
	var result []LocationRec
	var url string
	if strings.Contains(filesetId, `usx`) {
		url = HOST + "bibles/filesets/" + filesetId + "/ALL/1?v=4&limit=100000"
	} else {
		url = HOST + "download/" + filesetId + "?v=4"
	}
	content, status := d.httpGet(url, filesetId)
	if len(content) == 0 {
		return result
	}
	var response LocationDownloadRec
	err := json.Unmarshal(content, &response)
	if err != nil {
		log.Fatalln("Error parsing json for", filesetId, status, err)
	}
	return response.Data
}

func (d *DBPAPIClient) sortFileLocations(locations []LocationRec) []LocationRec {
	for i, loc := range locations {
		url, err := url.Parse(loc.URL)
		if err != nil {
			log.Fatalln("Could not parse URL", loc.URL, err)
		}
		locations[i].Filename = filepath.Base(url.Path)
	}
	sort.Slice(locations, func(i int, j int) bool {
		return locations[i].Filename < locations[j].Filename
	})
	return locations
}

func (d *DBPAPIClient) downloadFiles(directory string, locations []LocationRec) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		os.MkdirAll(directory, 0755)
	}
	for _, loc := range locations {
		size := loc.FileSize
		filePath := filepath.Join(directory, loc.Filename)
		file, err := os.Stat(filePath)
		if os.IsNotExist(err) || file.Size() != int64(size) {
			fmt.Println("Downloading", loc.Filename)
			content, status := d.httpGet(loc.URL, loc.Filename)
			if len(content) != size {
				log.Println("Warning for", loc.Filename, "has an expected size of", size, "but, actual size is", len(content))
			}
			if len(content) > 0 {
				d.saveFile(filePath, content)
			} else {
				log.Println("Error HTTP status", status)
			}
		}
	}
}

func (d *DBPAPIClient) httpGet(url string, desc string) ([]byte, string) {
	var body []byte
	var status string
	url += `&limit=100000&key=` + os.Getenv(`FCBH_DBP_KEY`)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error in DBP API request for:", desc, err)
		return body, status
	}
	defer resp.Body.Close()
	status = resp.Status
	if status[0] == '2' {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading DBP API response for:", desc, err)
			return body, status
		}
	}
	return body, status
}

func (d *DBPAPIClient) saveFile(filePath string, content []byte) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Error Creating file for download", err)
	}
	fp.Write(content)
	err = fp.Close()
	if err != nil {
		log.Fatalln("Error closing file for download", err)
	}
}
