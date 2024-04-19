package input

import "path/filepath"

type InputFile struct {
	MediaId    string
	MediaType  string
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
