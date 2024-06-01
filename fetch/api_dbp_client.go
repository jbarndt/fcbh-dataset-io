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
	BibleId            string       `json:"abbr"`
	LanguageISO        string       `json:"iso"`
	LanguageId         int          `json:"language_id"`
	RolvId             int          `json:"language_rolv_code"`
	LanguageName       string       `json:"language"`
	VersionName        string       `json:"name"`
	Alphabet           AlphabetType `json:"alphabet"` // alphabet.script
	Copyright          string       `json:"mark"`
	Books              []BooksType  `json:"books"`
	DbpProd            DbpProdType  `json:"filesets"`
	VersionCode        string
	AudioOTFileset     FilesetType
	AudioNTFileset     FilesetType
	TextOTPlainFileset FilesetType
	TextNTPlainFileset FilesetType
	TextOTUSXFileset   FilesetType
	TextNTUSXFileset   FilesetType
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
	textType := text.TextType()
	codec, bitrate := audio.AudioType()
	if testament.OT || len(testament.OTBooks) > 0 {
		info.TextOTPlainFileset = d.searchPlainText(info, `OT`, textType)
		info.TextOTUSXFileset = d.searchUSXText(info, `OT`, textType)
		info.AudioOTFileset = d.searchAudio(info, `OT`, `audio_drama`, codec, bitrate)
		tmpAudio := d.searchAudio(info, `OT`, `audio`, codec, bitrate)
		if tmpAudio.Id != `` {
			info.AudioOTFileset = tmpAudio
		}
	}
	if testament.NT || len(testament.NTBooks) > 0 {
		info.TextNTPlainFileset = d.searchPlainText(info, `NT`, textType)
		info.TextNTUSXFileset = d.searchUSXText(info, `NT`, textType)
		info.AudioNTFileset = d.searchAudio(info, `NT`, `audio_drama`, codec, bitrate)
		tmpAudio := d.searchAudio(info, `NT`, `audio`, codec, bitrate)
		if tmpAudio.Id != `` {
			info.AudioNTFileset = tmpAudio
		}
	}
}

func (d *APIDBPClient) searchPlainText(info *BibleInfoType, size string, textType request.MediaType) FilesetType {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == `text_plain` && rec.Size == size {
			if textType == request.TextPlain || textType == request.TextPlainEdit {
				return rec
			}
		}
	}
	return FilesetType{}
}

func (d *APIDBPClient) searchUSXText(info *BibleInfoType, size string, textType request.MediaType) FilesetType {
	for _, rec := range info.DbpProd.Filesets {
		if rec.Type == `text_usx` && rec.Size == size {
			if textType == request.TextUSXEdit || textType == request.TextPlainEdit {
				return rec
			}
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

func (d *APIDBPClient) UpdateIdent(id db.Ident, info BibleInfoType, textType request.MediaType) db.Ident {
	id.BibleId = info.BibleId
	if info.AudioOTFileset.Id != `` {
		id.AudioOTId = info.AudioOTFileset.Id
	}
	if info.AudioNTFileset.Id != `` {
		id.AudioNTId = info.AudioNTFileset.Id
	}
	if textType == request.TextPlain || textType == request.TextPlainEdit {
		if info.TextOTPlainFileset.Id != `` {
			id.TextOTId = info.TextOTPlainFileset.Id
		}
		if info.TextNTPlainFileset.Id != `` {
			id.TextNTId = info.TextNTPlainFileset.Id
		}
	} else if textType == request.TextUSXEdit {
		if info.TextOTUSXFileset.Id != `` {
			id.TextOTId = info.TextOTUSXFileset.Id
		}
		if info.TextNTUSXFileset.Id != `` {
			id.TextNTId = info.TextNTUSXFileset.Id
		}
	}
	if info.LanguageISO != `` {
		id.LanguageISO = info.LanguageISO
	}
	if info.VersionCode != `` {
		id.VersionCode = info.VersionCode
	}
	if info.LanguageId != 0 {
		id.LanguageId = info.LanguageId
	}
	if info.RolvId != 0 {
		id.RolvId = info.RolvId
	}
	if info.Alphabet.Alphabet != `` {
		id.Alphabet = info.Alphabet.Alphabet
	}
	if info.LanguageName != `` {
		id.LanguageName = info.LanguageName
	}
	if info.VersionName != `` {
		id.VersionName = info.VersionName
	}
	if textType != request.TextNone {
		id.TextSource = textType
	}
	return id
}
