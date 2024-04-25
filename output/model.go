package output

type Meta struct {
	Index int
	Name  string
	Tag   string
	Dtype string
	Cols  int // I don't
}

type Script struct {
	ScriptId      int     `name:"script_id,int"`
	BookId        string  `name:"book_id,string"`
	ChapterNum    int     `name:"chapter_num,int"`
	ChapterEnd    int     `name:"chapter_end,int"`
	AudioFile     string  `name:"audio_file,string"`
	ScriptNum     string  `name:"script_num,string"`
	UsfmStyle     string  `name:"usfm_style,string"`
	Person        string  `name:"person,string"`
	Actor         string  `name:"actor,string"`
	VerseStr      string  `name:"verse_str,string"`
	VerseEnd      string  `name:"verse_end,string"`
	ScriptText    string  `name:"script_text,string"`
	ScriptBeginTS float64 `name:"script_begin_ts,float64"`
	ScriptEndTS   float64 `name:"script_end_ts,float64"`
	MFCCRows      int
	MFCCCols      int
	MFCC          [][]float32 `name:"mfcc,[][]float32"`
}

type Word struct {
	WordId      int       `name:"word_id,int"`
	ScriptId    int       `name:"script_id,int"`
	BookId      string    `name:"book_id,string"`
	ChapterNum  int       `name:"chapter_num,int"`
	VerseStr    string    `name:"verse_str,string"`
	VerseNum    int       `name:"verse_num,int"`
	WordSeq     int       `name:"word_seq,int"`
	Word        string    `name:"word,string"`
	WordBeginTS float64   `name:"word_begin_ts,float64"`
	WordEndTS   float64   `name:"word_end_ts,float64"`
	WordEncoded []float64 `name:"word_enc,float64"`
	MFCCRows    int
	MFCCCols    int
	MFCC        [][]float32 `name:"mfcc,[][]float32"`
}

type HasMFCC interface {
	GetMFCC() [][]float32
	Rows() int
	Cols() int
	SetMFCC(mfcc [][]float32)
}

func (s *Script) GetMFCC() [][]float32 {
	return s.MFCC
}

func (s *Script) Rows() int {
	return s.MFCCRows
}

func (s *Script) Cols() int {
	return s.MFCCCols
}

func (s *Script) SetMFCC(mfcc [][]float32) {
	s.MFCC = mfcc
}
