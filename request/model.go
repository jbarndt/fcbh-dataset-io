package request

type Request struct {
	Required      Required      `yaml:"Required"`
	Testament     Testament     `yaml:"Testament,omitempty"`
	AudioData     AudioData     `yaml:"AudioData,omitempty"`
	TextData      TextData      `yaml:"TextData,omitempty"`
	Detail        Detail        `yaml:"Detail,omitempty"`
	Timestamps    Timestamps    `yaml:"Timestamps,omitempty"`
	AudioEncoding AudioEncoding `yaml:"AudioEncoding,omitempty"`
	TextEncoding  TextEncoding  `yaml:"TextEncoding,omitempty"`
	OutputFormat  OutputFormat  `yaml:"OutputFormat,omitempty"`
	Compare       Compare       `yaml:"Compare,omitempty"`
}

type Required struct {
	RequestName    string `yaml:"RequestName"`
	RequestorName  string `yaml:"RequestorName"`
	RequestorEmail string `yaml:"RequestorEmail"`
	BibleId        string `yaml:"BibleId"`
	LanguageISO    string `yaml:"LanguageISO"`
	VersionCode    string `yaml:"VersionCode"`
}

type Testament struct {
	NT      bool     `yaml:"NT,omitempty"`
	NTBooks []string `yaml:"NTBooks,omitempty"`
	OT      bool     `yaml:"OT,omitempty"`
	OTBooks []string `yaml:"OTBooks,omitempty"`
}

/*
	func (t Testament) String() string {
		var result string
		if t.NT && t.OT {
			result = `C`
		}
		if t.NT {
			result = `NT`
		}
		if t.OT {
			result = `OT`
		}
		return result
	}
*/
type AudioData struct {
	BibleBrain BibleBrainAudio `yaml:"BibleBrain,omitempty"`
	File       string          `yaml:"File,omitempty"`
	Http       string          `yaml:"Http,omitempty"`
	AWSS3      string          `yaml:"AWSS3,omitempty"`
	POST       bool            `yaml:"POST,omitempty"`
	NoAudio    bool            `yaml:"NoAudio,omitempty"`
}

type BibleBrainAudio struct {
	MP3_64 bool `yaml:"MP3_64,omitempty"`
	MP3_16 bool `yaml:"MP3_16,omitempty"`
	OPUS   bool `yaml:"OPUS,omitempty"`
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
	BibleBrain   BibleBrainText `yaml:"BibleBrain,omitempty"`
	SpeechToText SpeechToText   `yaml:"SpeechToText,omitempty"`
	File         string         `yaml:"File,omitempty"`
	Http         string         `yaml:"Http,omitempty"`
	AWSS3        string         `yaml:"AWSS3,omitempty"`
	POST         bool           `yaml:"POST,omitempty"`
	NoText       bool           `yaml:"NoText,omitempty"`
}

type BibleBrainText struct {
	TextUSXEdit   bool `yaml:"TextUSXEdit,omitempty"`
	TextPlainEdit bool `yaml:"TextPlainEdit,omitempty"`
	TextPlain     bool `yaml:"TextPlain,omitempty"`
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
	Whisper Whisper `yaml:"Whisper,omitempty"`
}

type Whisper struct {
	Model WhisperModel `yaml:"Model,omitempty"`
}
type WhisperModel struct {
	Large  bool `yaml:"Large,omitempty"`
	Medium bool `yaml:"Medium,omitempty"`
	Small  bool `yaml:"Small,omitempty"`
	Base   bool `yaml:"Base,omitempty"`
	Tiny   bool `yaml:"Tiny,omitempty"`
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
	Lines bool `yaml:"Lines,omitempty"`
	Words bool `yaml:"Words,omitempty"`
}

type Timestamps struct {
	BibleBrain   bool `yaml:"BibleBrain,omitempty"`
	Aeneas       bool `yaml:"Aeneas,omitempty"`
	NoTimestamps bool `yaml:"NoTimestamps,omitempty"`
}

type AudioEncoding struct {
	MFCC       bool `yaml:"MFCC,omitempty"`
	NoEncoding bool `yaml:"NoEncoding,omitempty"`
}

type TextEncoding struct {
	FastText   bool `yaml:"FastText,omitempty"`
	NoEncoding bool `yaml:"NoEncoding,omitempty"`
}

type OutputFormat struct {
	CSV    bool `yaml:"CSV,omitempty"`
	JSON   bool `yaml:"JSON,omitempty"`
	Sqlite bool `yaml:"Sqlite,omitempty"`
}

type Compare struct {
	Project1        string          `yaml:"Project1,omitempty"`
	Project2        string          `yaml:"Project2,omitempty"`
	CompareSettings CompareSettings `yaml:"CompareSettings,omitempty"`
}

type CompareSettings struct {
	LowerCase         bool              `yaml:"LowerCase,omitempty"`
	RemovePromptChars bool              `yaml:"RemovePromptChars,omitempty"`
	RemovePunctuation bool              `yaml:"RemovePunctuation,omitempty"`
	DoubleQuotes      CompareChoice     `yaml:"DoubleQuotes,omitempty"`
	Apostrophe        CompareChoice     `yaml:"Apostrophe,omitempty"`
	Hyphen            CompareChoice     `yaml:"Hyphen,omitempty"`
	DiacriticalMarks  DiacriticalChoice `yaml:"DiacriticalMarks,omitempty"`
}

type CompareChoice struct {
	Remove    bool `yaml:"Remove,omitempty"`
	Normalize bool `yaml:"Normalize,omitempty"`
}

type DiacriticalChoice struct {
	Remove        bool `yaml:"Remove,omitempty"`
	NormalizeNFC  bool `yaml:"NormalizeNFC,omitempty"`
	NormalizeNFD  bool `yaml:"NormalizeNFD,omitempty"`
	NormalizeNFKC bool `yaml:"NormalizeNFKC,omitempty"`
	NormalizeNFKD bool `yaml:"NormalizeNFKD,omitempty"`
}
