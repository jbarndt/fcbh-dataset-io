package req

type Request struct {
	Required      Required      `yaml:"Required"`
	Testament     Testament     `yaml:"Testament"`
	AudioData     AudioData     `yaml:"AudioData"`
	TextData      TextData      `yaml:"TextData"`
	Detail        Detail        `yaml:"Detail"`
	Timestamps    Timestamps    `yaml:"Timestamps"`
	AudioEncoding AudioEncoding `yaml:"AudioEncoding"`
	TextEncoding  TextEncoding  `yaml:"TextEncoding"`
	OutputFormat  OutputFormat  `yaml:"OutputFormat"`
	Compare       Compare       `yaml:"Compare"`
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
	NT bool `yaml:"NT"`
	OT bool `yaml:"OT"`
}

type AudioData struct {
	BibleBrain BibleBrainAudio `yaml:"BibleBrain"`
	File       string          `yaml:"File"`
	Http       string          `yaml:"Http"`
	AWSS3      string          `yaml:"AWSS3"`
	POST       bool            `yaml:"POST"`
	NoAudio    bool            `yaml:"NoAudio"`
}

type BibleBrainAudio struct {
	MP3_64 bool `yaml:"MP3_64"`
	MP3_16 bool `yaml:"MP3_16"`
	OPUS   bool `yaml:"OPUS"`
}

type TextData struct {
	BibleBrain   BibleBrainText `yaml:"BibleBrain"`
	SpeechToText SpeechToText   `yaml:"SpeechToText"`
	File         string         `yaml:"File"`
	Http         string         `yaml:"Http"`
	AWSS3        string         `yaml:"AWSS3"`
	POST         bool           `yaml:"POST"`
	NoText       bool           `yaml:"NoText"`
}

type BibleBrainText struct {
	TextUSXEdit   bool `yaml:"TextUSXEdit"`
	TextPlainEdit bool `yaml:"TextPlainEdit"`
	TextPlain     bool `yaml:"TextPlain"`
}

type SpeechToText struct {
	Whisper Whisper `yaml:"Whisper"`
}

type Whisper struct {
	Model WhisperModel `yaml:"Model"`
}
type WhisperModel struct {
	Large  bool `yaml:"Large"`
	Medium bool `yaml:"Medium"`
	Small  bool `yaml:"Small"`
	Base   bool `yaml:"Base"`
	Tiny   bool `yaml:"Tiny"`
}

type Detail struct {
	Lines bool `yaml:"Lines"`
	Words bool `yaml:"Words"`
}

type Timestamps struct {
	BibleBrain   bool `yaml:"BibleBrain"`
	Aeneas       bool `yaml:"Aeneas"`
	NoTimestamps bool `yaml:"NoTimestamps"`
}

type AudioEncoding struct {
	MFCC       bool `yaml:"MFCC"`
	NoEncoding bool `yaml:"NoEncoding"`
}

type TextEncoding struct {
	FastText   bool `yaml:"FastText"`
	NoEncoding bool `yaml:"NoEncoding"`
}

type OutputFormat struct {
	CSV    bool `yaml:"CSV"`
	JSON   bool `yaml:"JSON"`
	Sqlite bool `yaml:"Sqlite"`
}

type Compare struct {
	Project1        string          `yaml:"Project1"`
	Project2        string          `yaml:"Project2"`
	CompareSettings CompareSettings `yaml:"CompareSettings"`
}

type CompareSettings struct {
	LowerCase         bool              `yaml:"LowerCase"`
	RemovePromptChars bool              `yaml:"RemovePromptChars"`
	RemovePunctuation bool              `yaml:"RemovePunctuation"`
	DoubleQuotes      CompareChoice     `yaml:"DoubleQuotes"`
	Apostrophe        CompareChoice     `yaml:"Apostrophe"`
	Hyphen            CompareChoice     `yaml:"Hyphen"`
	DiacriticalMarks  DiacriticalChoice `yaml:"DiacriticalMarks"`
}

type CompareChoice struct {
	Remove    bool `yaml:"Remove"`
	Normalize bool `yaml:"Normalize"`
}

type DiacriticalChoice struct {
	Remove        bool `yaml:"Remove"`
	NormalizeNFC  bool `yaml:"NormalizeNFC"`
	NormalizeNFD  bool `yaml:"NormalizeNFD"`
	NormalizeNFKC bool `yaml:"NormalizeNFKC"`
	NormalizeNFKD bool `yaml:"NormalizeNFKD"`
}
