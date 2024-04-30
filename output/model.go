package output

type Meta struct {
	Index  int // Position of field in the struct
	Name   string
	Tag    string
	Dtype  string // Only used by json_writer to control double quotes
	CSVPos int    // Position of field in CSV output
	Cols   int
}

type Script struct {
	ScriptId      int     `name:"script_id,int"`
	Reference     string  `name:"reference,string"`
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
	MFCC          [][]float64 `name:"mfcc,[][]float64"`
}

type Word struct {
	WordId      int       `name:"word_id,int"`
	ScriptId    int       `name:"script_id,int"`
	Reference   string    `name:"reference,string"`
	BookId      string    `name:"book_id,string"`
	ChapterNum  int       `name:"chapter_num,int"`
	ChapterEnd  int       `name:"chapter_end,int"`
	VerseStr    string    `name:"verse_str,string"`
	VerseEnd    string    `name:"verse_end,string"`
	VerseNum    int       `name:"verse_num,int"`
	UsfmStyle   string    `name:"usfm_style,string"`
	Person      string    `name:"person,string"`
	Actor       string    `name:"actor,string"`
	WordSeq     int       `name:"word_seq,int"`
	Word        string    `name:"word,string"`
	WordBeginTS float64   `name:"word_begin_ts,float64"`
	WordEndTS   float64   `name:"word_end_ts,float64"`
	WordEnc     []float64 `name:"word_enc,float64"`
	MFCCRows    int
	MFCCCols    int
	MFCC        [][]float64 `name:"mfcc,[][]float64"`
}
