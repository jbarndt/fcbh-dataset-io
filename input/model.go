package input

import (
	"dataset/decode_yaml/request"
	"path/filepath"
)

type InputFile struct {
	MediaId    string
	MediaType  request.MediaType
	Testament  string
	BookId     string // not used for text_plain
	BookSeq    string
	Chapter    int    // only used for audio
	Verse      string // used by OBT and Vessel
	ChapterEnd int
	VerseEnd   string
	Filename   string
	FileExt    string
	Directory  string
}

func (i *InputFile) FilePath() string {
	return filepath.Join(i.Directory, i.Filename)
}
