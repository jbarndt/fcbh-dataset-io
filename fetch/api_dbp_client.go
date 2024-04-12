package fetch

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"dataset/request"
	"encoding/json"
	"strings"
)

type APIDBPClient struct {
	ctx     context.Context
	bibleId string
}

func NewAPIDBPClient(ctx context.Context, bibleId string) APIDBPClient {
	var d APIDBPClient
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

func (d *APIDBPClient) BibleInfo() (BibleInfoType, dataset.Status) {
	var result BibleInfoType
	var status dataset.Status
	var get = `https://4.dbt.io/api/bibles/` + d.bibleId + `?v=4`
	var response BibleInfoRespType
	body, status := httpGet(d.ctx, get, d.bibleId)
	if !status.IsErr {
		err := json.Unmarshal(body, &response)
		if err != nil {
			status := log.Error(d.ctx, 500, err, "Error decoding DBP API /bibles JSON")
			return result, status
		}
		result = response.Data
		result.VersionCode = d.bibleId[3:]
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

func (d *APIDBPClient) FindFilesets(info *BibleInfoType, audio request.BibleBrainAudio,
	text request.BibleBrainText, testament request.Testament) bool {
	var okAudio = true
	var okText = true
	var reqSize = testament.String()
	codec, bitrate := audio.AudioType()
	if codec != `` {
		okAudio = d.searchAudio(info, `audio`, reqSize, codec, bitrate)
		if !okAudio {
			okAudio = d.searchAudio(info, `audio_drama`, reqSize, codec, bitrate)
			if !okAudio {
				if codec == `MP3` && bitrate == `64kbps` {
					okAudio = d.searchAudio(info, `audio`, reqSize, `MP`, `3kbps`)
					if !okAudio {
						okAudio = d.searchAudio(info, `audio_drama`, reqSize, `MP`, `3kbps`)
					}
				}
			}
		}
	}
	if text.String() != `` {
		okText = d.searchText(info, text.String(), reqSize)
	}
	return okAudio && okText
}

func (d *APIDBPClient) searchText(info *BibleInfoType, reqType string, reqSize string) bool {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == reqType && (reqSize == `C` || rec.Size == reqSize) {
			info.TextFilesets = append(info.TextFilesets, rec)
		}
	}
	return len(info.TextFilesets) > 0
}

func (d *APIDBPClient) searchAudio(info *BibleInfoType, reqType string, reqSize string, codec string, bitRate string) bool {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == reqType && (reqSize == `C` || rec.Size == reqSize) {
			if strings.ToUpper(rec.Codec) == codec || strings.ToUpper(rec.Codec)+`3` == codec {
				if rec.Bitrate == bitRate || bitRate == `` {
					info.AudioFilesets = append(info.AudioFilesets, rec)
				}
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
