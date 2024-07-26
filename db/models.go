package db

import "dataset/request"

type Ident struct {
	DatasetId    int64
	BibleId      string
	AudioOTId    string
	AudioNTId    string
	TextOTId     string
	TextNTId     string
	TextSource   request.MediaType
	LanguageISO  string
	VersionCode  string
	LanguageId   int
	RolvId       int
	Alphabet     string
	LanguageName string
	VersionName  string
}

type Script struct {
	ScriptId      int
	DatasetId     int
	BookId        string
	ChapterNum    int
	ChapterEnd    int
	AudioFile     string
	ScriptNum     string
	UsfmStyle     string
	Person        string
	Actor         string
	VerseNum      int
	VerseStr      string
	VerseEnd      string
	ScriptText    string
	ScriptTexts   []string
	ScriptBeginTS float64
	ScriptEndTS   float64
}

type Word struct {
	WordId      int
	ScriptId    int
	WordSeq     int
	VerseNum    int
	TType       string
	Word        string
	WordBeginTS float64
	WordEndTS   float64
	WordEncoded []float64
}

type Timestamp struct {
	Id        int
	VerseStr  string
	AudioFile string
	Text      string
	BeginTS   float64
	EndTS     float64
}

type MFCC struct {
	Id   int
	Rows int
	Cols int
	MFCC [][]float32
}

type Language struct {
	GlottoId    string `json:"id"`
	FamilyId    string `json:"family_id"`
	ParentId    string `json:"parent_id"`
	Name        string `json:"name"`
	Bookkeeping bool   `json:"bookkeeping"`
	Level       string `json:"level"` //(language, dialect, family)
	Iso6393     string `json:"iso639_3"`
	CountryIds  string `json:"country_ids"`
	Iso6391     string `json:"iso639_1"`
	Whisper     bool   `json:"whisper"`
	MMSASR      bool   `json:"mms_asr"`
	ESpeak      bool   `json:"espeak"`
}
