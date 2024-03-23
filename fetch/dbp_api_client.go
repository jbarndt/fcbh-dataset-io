package fetch

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type DbpApiClient struct {
	bibleId string
	//audioSource dataset_io.AudioSourceType
}

func NewDbpApiClient(bibleId string) *DbpApiClient {
	var d DbpApiClient
	d.bibleId = bibleId
	return &d
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
	DbpProd []FilesetType `json:"dbp-prod"`
}
type AlphabetType struct {
	Alphabet string `json:"script"`
}
type BibleInfoType struct {
	BibleId      string `json:"abbr"`
	LanguageISO  string `json:"iso"`
	VersionCode  string
	LanguageId   int          `json:"language_id"`
	RolvId       int          `json:"language_rolv_code"`
	LanguageName string       `json:"language"`
	VersionName  string       `json:"name"`
	Alphabet     AlphabetType `json:"alphabet"` // alphabet.script
	DbpProd      DbpProdType  `json:"filesets"`
}
type BibleInfoRespType struct {
	Data BibleInfoType `json:"data""`
}

func (d *DbpApiClient) BibleInfo() BibleInfoType {
	var url = `https://4.dbt.io/api/bibles/ATIWBT?`
	var response BibleInfoRespType
	body := d.query(url)
	if body != nil && len(body) > 0 {
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Println("Error decoding DBP API JSON:", err)
			return BibleInfoType{}
		}
	}
	var result = response.Data
	result.VersionCode = d.bibleId[3:]
	return result
}

func (d *DbpApiClient) query(url string) []byte {
	url += `v=4&limit=100000&key=` + os.Getenv(`FCBH_DBP_KEY`)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error in DBP API request:", err)
		return []byte{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading DBP API response:", err)
		return []byte{}
	}
	return body
}
