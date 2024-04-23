package output

type Script struct {
	ScriptId      int    `name:"script_id,int"`
	DatasetId     int    `name:"dataset_id,int"`
	BookId        string `name:"book_id,string"`
	ChapterNum    int    `name:"chapter_num,int"`
	ChapterEnd    int    `name:"chapter_end,int"`
	AudioFile     string `name:"audio_file,string"`
	ScriptNum     string `name:"script_num,string"`
	UsfmStyle     string `name:"usfm_style,string"`
	Person        string `name:"person,string"`
	Actor         string `name:"actor,string"`
	VerseNum      int    `name:"verse_num,int"`
	VerseStr      string `name:"verse_str,string"`
	VerseEnd      string `name:"verse_end,string"`
	ScriptText    string `name:"script_text,string"`
	ScriptTexts   []string
	ScriptBeginTS float64 `name:"script_begin_ts,float64"`
	ScriptEndTS   float64 `name:"script_end_ts,float64"`
	MFCC          MFCC
}

type MFCC struct {
	Id   int         `name:"id,int"`
	Rows int         `name:"rows,int"`
	Cols int         `name:"cols,int"`
	MFCC [][]float32 `name:"mfcc,[][]float32"`
}
