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

type Audio struct {
	Book         string  `json:"book"`
	ChapterNum   int     `json:"chapter_num"`
	ChapterEnd   int     `json:"chapter_end"`
	VerseStr     string  `json:"verse_str"`
	VerseEnd     string  `json:"verse_end"`
	VerseSeq     int     `json:"verse_seq"`
	BeginTS      float64 `json:"begin_ts"`
	EndTS        float64 `json:"end_ts"`
	Score        float64 `json:"score"`
	Uroman       string  `json:"uroman"` // Is this needed
	Text         string  `json:"text"`
	AudioChapter string  `json:"audio_chapter"`
	AudioVerse   string  `json:"-"` // This exists temporarily
}
