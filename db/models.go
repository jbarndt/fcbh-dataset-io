package db

type Ident struct {
	DatasetId    int
	BibleId      string
	AudioOTId    string
	AudioNTId    string
	TextOTId     string
	TextNTId     string
	TextSource   string
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
	AudioFile     string
	ScriptNum     string
	UsfmStyle     string
	Person        string
	Actor         string
	VerseNum      int
	VerseStr      string
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
