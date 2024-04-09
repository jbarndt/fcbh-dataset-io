package dataset

type SourceType string

const (
	API SourceType = "API"
	CLI SourceType = "CLI"
)

type TestamentType string

const (
	NT                   TestamentType = "NT"
	OT                   TestamentType = "OT"
	C                    TestamentType = "C"
	DefaultTestamentType               = C
)

type DetailType string

const (
	LINES                 DetailType = "LINES"
	WORDS                 DetailType = "WORDS"
	BOTH                  DetailType = "BOTH"
	DefaultTextDetailType            = BOTH
)

type AudioSourceType string

const (
	MP3_64                 AudioSourceType = "MP3_64"
	MP3_16                 AudioSourceType = "MP3_16"
	OPUS                   AudioSourceType = "OPUS"
	NOAUDIO                AudioSourceType = "NOAUDIO"
	DefaultAudioSourceType                 = MP3_64
)

type TextSourceType string

const (
	SCRIPT   TextSourceType = "SCRIPT"
	DBPTEXT  TextSourceType = "DBPTEXT"
	TEXTEDIT TextSourceType = "TEXTEDIT"
	USXEDIT  TextSourceType = "USXEDIT"
	NOTEXT   TextSourceType = "NOTEXT"
	//SPEECHTOTEXT          TextSourceType = "SPEECHTOTEXT"
	DefaultTextSourceType = USXEDIT
)

type SpeechToTextType string

const (
	WHISPER             SpeechToTextType = "WHISPER"
	NOSTT               SpeechToTextType = "NOSTT"
	DefaultSpeechToText                  = NOSTT
)

type AudioEncodingType string

const (
	TIMESTAMP_ONLY           AudioEncodingType = "TIMESTAMP_ONLY"
	MFCC                     AudioEncodingType = "MFCC"
	NOAUDIO_ENCODE           AudioEncodingType = "NOAUDIO_ENCODE"
	DefaultAudioEncodingType                   = NOAUDIO_ENCODE
)

type TextEncodingType string

const (
	FASTTEXT                TextEncodingType = "FASTTEXT"
	DefaultTextEncodingType                  = FASTTEXT
)

type OutputFormatType string

const (
	JSON                OutputFormatType = "JSON"
	CSV                 OutputFormatType = "CSV"
	SQLITE              OutputFormatType = "SQLITE"
	DefaultOutputFormat                  = JSON
)

type RequestType struct {
	Email          string
	BibleId        string
	AudioFilesetId string
	TextFilesetId  string
	Source         SourceType
	Testament      TestamentType
	AudioSource    AudioSourceType
	TextDetail     DetailType
	TextSource     TextSourceType
	SpeechToText   SpeechToTextType
	AudioEncoding  AudioEncodingType
	TextEncoding   TextEncodingType
	OutputForm     OutputFormatType
}

// prerquisites
// fetch before read, can be in the code
// although there can be an explicit fetch
