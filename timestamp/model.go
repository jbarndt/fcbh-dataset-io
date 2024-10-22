package timestamp

type Timestamp struct {
	Book         string  `json:"book"`
	Chapter      int     `json:"chapter"`
	Verse        string  `json:"verse"`
	Seq          int     `json:"seq"`
	BeginTS      float64 `json:"begin_ts"`
	EndTS        float64 `json:"end_ts"`
	Score        float64 `json:"score"`
	Uroman       string  `json:"uroman"` // Is this needed
	Text         string  `json:"text"`
	AudioChapter string  `json:"audio_chapter"`
	AudioVerse   string  `json:"-"` // This exists temporarily
}

// There could be a media type with mediaId, lang, and hold books
// There could be a book type with bookId, and hold chapter
// There could be a chapter type that holds chapter_num, AudioChapter, and holds verses
// There could be a verse type that contains all of the rest

// In mms_asr we really only need the chapter and verse sections

// Except, when we update the database
