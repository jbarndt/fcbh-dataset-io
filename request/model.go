package request

type Request struct {
	IsNew         bool          `yaml:"isnew"`
	RequestName   string        `yaml:"requestname"`
	BibleId       string        `yaml:"bibleid"`
	Testament     Testament     `yaml:"testament,omitempty"`
	AudioData     AudioData     `yaml:"audiodata,omitempty"`
	TextData      TextData      `yaml:"textdata,omitempty"`
	Detail        Detail        `yaml:"detail,omitempty"`
	Timestamps    Timestamps    `yaml:"timestamps,omitempty"`
	AudioEncoding AudioEncoding `yaml:"audioencoding,omitempty"`
	TextEncoding  TextEncoding  `yaml:"textencoding,omitempty"`
	OutputFormat  OutputFormat  `yaml:"outputformat,omitempty"`
	Compare       Compare       `yaml:"compare,omitempty"`
}

type Testament struct {
	NT      bool     `yaml:"nt,omitempty"`
	NTBooks []string `yaml:"ntbooks,omitempty"`
	OT      bool     `yaml:"ot,omitempty"`
	OTBooks []string `yaml:"otbooks,omitempty"`
	otMap   map[string]bool
	ntMap   map[string]bool
}

func (t *Testament) BuildBookMaps() {
	t.otMap = make(map[string]bool)
	for _, book := range t.OTBooks {
		t.otMap[book] = true
	}
	t.ntMap = make(map[string]bool)
	for _, book := range t.NTBooks {
		t.ntMap[book] = true
	}
}

func (t *Testament) Has(ttype string, bookId string) bool {
	if ttype == `NT` {
		return t.HasNT(bookId)
	} else {
		return t.HasOT(bookId)
	}
}

func (t *Testament) HasOT(bookId string) bool {
	if t.OT {
		return true
	}
	_, ok := t.otMap[bookId]
	return ok
}

func (t *Testament) HasNT(bookId string) bool {
	if t.NT {
		return true
	}
	_, ok := t.ntMap[bookId]
	return ok
}

type AudioData struct {
	BibleBrain BibleBrainAudio `yaml:"biblebrain,omitempty"`
	File       string          `yaml:"file,omitempty"`
	Http       string          `yaml:"http,omitempty"`
	AWSS3      string          `yaml:"awss3,omitempty"`
	POST       bool            `yaml:"post,omitempty"`
	NoAudio    bool            `yaml:"noaudio,omitempty"`
}

type BibleBrainAudio struct {
	MP3_64 bool `yaml:"mp3_64,omitempty"`
	MP3_16 bool `yaml:"mp3_16,omitempty"`
	OPUS   bool `yaml:"opus,omitempty"`
}

func (b BibleBrainAudio) AudioType() (string, string) {
	var result string
	var codec string
	if b.MP3_64 {
		result = `MP3`
		codec = `64kbps`
	} else if b.MP3_16 {
		result = `MP3`
		codec = `16kbps`
	} else if b.OPUS {
		result = `OPUS`
		codec = ``
	}
	return result, codec
}

type TextData struct {
	BibleBrain   BibleBrainText `yaml:"biblebrain,omitempty"`
	SpeechToText SpeechToText   `yaml:"speechtotext,omitempty"`
	File         string         `yaml:"file,omitempty"`
	Http         string         `yaml:"http,omitempty"`
	AWSS3        string         `yaml:"awss3,omitempty"`
	POST         bool           `yaml:"post,omitempty"`
	NoText       bool           `yaml:"notext,omitempty"`
}

type BibleBrainText struct {
	TextUSXEdit   bool `yaml:"textusxedit,omitempty"`
	TextPlainEdit bool `yaml:"textplainedit,omitempty"`
	TextPlain     bool `yaml:"textplain,omitempty"`
}

func (b BibleBrainText) String() string {
	var result string
	if b.TextUSXEdit {
		result = `text_usx`
	} else if b.TextPlainEdit {
		result = `text_plain`
	} else if b.TextPlain {
		result = `text_plain`
	}
	return result
}

type SpeechToText struct {
	Whisper Whisper `yaml:"whisper,omitempty"`
}

type Whisper struct {
	Model WhisperModel `yaml:"model,omitempty"`
}
type WhisperModel struct {
	Large  bool `yaml:"large,omitempty"`
	Medium bool `yaml:"medium,omitempty"`
	Small  bool `yaml:"small,omitempty"`
	Base   bool `yaml:"base,omitempty"`
	Tiny   bool `yaml:"tiny,omitempty"`
}

func (w WhisperModel) String() string {
	var result string
	if w.Large {
		result = `large`
	} else if w.Medium {
		result = `medium`
	} else if w.Small {
		result = `small`
	} else if w.Base {
		result = `base`
	} else if w.Tiny {
		result = `tiny`
	}
	return result
}

type Detail struct {
	Lines bool `yaml:"lines,omitempty"`
	Words bool `yaml:"words,omitempty"`
}

type Timestamps struct {
	BibleBrain   bool `yaml:"biblebrain,omitempty"`
	Aeneas       bool `yaml:"aeneas,omitempty"`
	NoTimestamps bool `yaml:"notimestamps,omitempty"`
}

type AudioEncoding struct {
	MFCC       bool `yaml:"mfcc,omitempty"`
	NoEncoding bool `yaml:"noencoding,omitempty"`
}

type TextEncoding struct {
	FastText   bool `yaml:"fasttext,omitempty"`
	NoEncoding bool `yaml:"noencoding,omitempty"`
}

type OutputFormat struct {
	CSV    bool `yaml:"csv,omitempty"`
	JSON   bool `yaml:"json,omitempty"`
	Sqlite bool `yaml:"sqlite,omitempty"`
}

type Compare struct {
	BaseProject     string          `yaml:"baseproject,omitempty"`
	CompareSettings CompareSettings `yaml:"comparesettings,omitempty"`
}

type CompareSettings struct {
	LowerCase         bool              `yaml:"lowercase,omitempty"`
	RemovePromptChars bool              `yaml:"removepromptchars,omitempty"`
	RemovePunctuation bool              `yaml:"removepunctuation,omitempty"`
	DoubleQuotes      CompareChoice     `yaml:"doublequotes,omitempty"`
	Apostrophe        CompareChoice     `yaml:"apostrophe,omitempty"`
	Hyphen            CompareChoice     `yaml:"hyphen,omitempty"`
	DiacriticalMarks  DiacriticalChoice `yaml:"diacriticalmarks,omitempty"`
}

type CompareChoice struct {
	Remove    bool `yaml:"remove,omitempty"`
	Normalize bool `yaml:"normalize,omitempty"`
}

type DiacriticalChoice struct {
	Remove        bool `yaml:"remove,omitempty"`
	NormalizeNFC  bool `yaml:"normalizenfc,omitempty"`
	NormalizeNFD  bool `yaml:"normalizenfd,omitempty"`
	NormalizeNFKC bool `yaml:"normalizenfkc,omitempty"`
	NormalizeNFKD bool `yaml:"normalizenfkd,omitempty"`
}
