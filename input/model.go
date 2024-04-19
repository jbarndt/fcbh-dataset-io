package input

type InputFile struct {
	MediaId    string
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
