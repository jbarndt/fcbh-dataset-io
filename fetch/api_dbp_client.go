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
type BooksType struct {
	BookId    string `json:"book_id"`
	Name      string `json:"name"`
	Chapters  []int  `json:"chapters"`
	BookSeq   string `json:"book_seq"`
	Testament string `json:"testament"`
}
type AlphabetType struct {
	Alphabet string `json:"script"`
}
type BibleInfoType struct {
	BibleId        string       `json:"abbr"`
	LanguageISO    string       `json:"iso"`
	LanguageId     int          `json:"language_id"`
	RolvId         int          `json:"language_rolv_code"`
	LanguageName   string       `json:"language"`
	VersionName    string       `json:"name"`
	Alphabet       AlphabetType `json:"alphabet"` // alphabet.script
	Copyright      string       `json:"mark"`
	Books          []BooksType  `json:"books"`
	DbpProd        DbpProdType  `json:"filesets"`
	VersionCode    string
	AudioOTFileset FilesetType
	AudioNTFileset FilesetType
	TextOTFileset  FilesetType
	TextNTFileset  FilesetType
}
type BibleInfoRespType struct {
	Data BibleInfoType `json:"data"`
}

func (d *APIDBPClient) BibleInfo() (BibleInfoType, dataset.Status) {
	var result BibleInfoType
	var status dataset.Status
	var get = `https://4.dbt.io/api/bibles/` + d.bibleId + `?v=4`
	var response BibleInfoRespType
	body, status := httpGet(d.ctx, get, false, d.bibleId)
	if status.IsErr {
		return result, status
	}
	//fmt.Println(string(body))
	err := json.Unmarshal(body, &response)
	if err != nil {
		status := log.Error(d.ctx, 500, err, "Error decoding DBP API /bibles JSON")
		return result, status
	}
	result = response.Data
	result.VersionCode = d.bibleId[3:]
	return result, status
}

func (d *APIDBPClient) FindFilesets(info *BibleInfoType, audio request.BibleBrainAudio,
	text request.BibleBrainText, testament request.Testament) {
	textType := text.String()
	codec, bitrate := audio.AudioType()
	if testament.OT || len(testament.OTBooks) > 0 {
		info.TextOTFileset = d.searchText(info, `OT`, textType)
		info.AudioOTFileset = d.searchAudio(info, `OT`, `audio_drama`, codec, bitrate)
		tmpAudio := d.searchAudio(info, `OT`, `audio`, codec, bitrate)
		if tmpAudio.Id != `` {
			info.AudioOTFileset = tmpAudio
		}
	}
	if testament.NT || len(testament.NTBooks) > 0 {
		info.TextNTFileset = d.searchText(info, `NT`, textType)
		info.AudioNTFileset = d.searchAudio(info, `NT`, `audio_drama`, codec, bitrate)
		tmpAudio := d.searchAudio(info, `NT`, `audio`, codec, bitrate)
		if tmpAudio.Id != `` {
			info.AudioNTFileset = tmpAudio
		}
	}
}

func (d *APIDBPClient) searchText(info *BibleInfoType, size string, textType string) FilesetType {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == textType && rec.Size == size {
			return rec
		}
	}
	return FilesetType{}
}

func (d *APIDBPClient) searchAudio(info *BibleInfoType, size string, audioType string, codec string, bitrate string) FilesetType {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == audioType {
			recCodec := strings.ToUpper(rec.Codec)
			if recCodec == codec || (recCodec == `MP` && codec == `MP3`) {
				if rec.Bitrate == bitrate || (rec.Bitrate == `3kbps` && bitrate == `64kbps`) {
					if rec.Size == size {
						if recCodec == `MP` {
							rec.Codec = `MP3`
						}
						if rec.Bitrate == `3kbps` {
							rec.Bitrate = `64kbps`
						}
						return rec
					}
				}
			}
		}
	}
	return FilesetType{}
}

func (d *APIDBPClient) CreateIdent(info BibleInfoType) db.Ident {
	var id db.Ident
	id.BibleId = info.BibleId
	id.AudioOTId = info.AudioOTFileset.Id
	id.AudioNTId = info.AudioNTFileset.Id
	id.TextOTId = info.TextOTFileset.Id
	id.TextNTId = info.TextNTFileset.Id
	id.LanguageISO = info.LanguageISO
	id.VersionCode = info.VersionCode
	id.LanguageId = info.LanguageId
	id.RolvId = info.RolvId
	id.Alphabet = info.Alphabet.Alphabet
	id.LanguageName = info.LanguageName
	id.VersionName = info.VersionName
	return id
}
