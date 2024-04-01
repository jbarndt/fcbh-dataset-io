package read

import (
	"context"
	"dataset"
	"dataset/db"
	"strconv"
)

/**
This class uses DBPTextReader to get the Script records
But then it gets the mt and h headings from USX files,
And add this to the DBPText data, and stores the combined data.
*/

//Could this all be done is SQL?
//Select from each source 1 book at a time, and one chapter at a time

type DBPTextEditReader struct {
	ctx     context.Context
	bibleId string
	conn    db.DBAdapter
}

func NewDBPTextEditReader(bibleId string, conn db.DBAdapter) DBPTextEditReader {
	var d DBPTextEditReader
	d.ctx = conn.Ctx
	d.bibleId = bibleId
	d.conn = conn
	return d
}

func (d *DBPTextEditReader) Process(testament dataset.TestamentType) {
	d.ensureDBPText(testament)
	d.ensureUSXEDITText(testament)
	titleMap, chapMap := d.readUSXHeadings(testament)
	records := d.combineHeadings(titleMap, chapMap)
	d.conn.InsertScripts(records)
}

func (d *DBPTextEditReader) ensureDBPText(testament dataset.TestamentType) {
	var sourceDB = d.bibleId + `_DBPTEXT.db`
	if !db.Exists(sourceDB) {
		var conn2 = db.NewDBAdapter(d.ctx, sourceDB)
		textAdapter := NewDBPTextReader(conn2)
		textAdapter.ProcessDirectory(d.bibleId, testament)
		conn2.Close()
	}
}

func (d *DBPTextEditReader) ensureUSXEDITText(testament dataset.TestamentType) {
	var sourceDB = d.bibleId + `_USXEDIT.db`
	if !db.Exists(sourceDB) {
		var conn3 = db.NewDBAdapter(d.ctx, sourceDB)
		ReadUSXEdit(conn3, d.bibleId, testament)
		conn3.Close()
	}
}

func (d *DBPTextEditReader) readUSXHeadings(testament dataset.TestamentType) (map[string][]db.Script, map[string]db.Script) {
	var sourceDB = d.bibleId + `_USXEDIT.db`
	var conn4 = db.NewDBAdapter(d.ctx, sourceDB)
	records := conn4.SelectScriptHeadings()
	conn4.Close()
	var bookTitle = make(map[string][]db.Script)
	var chapTitle = make(map[string]db.Script)
	for _, rec := range records {
		if rec.UsfmStyle == `para.h` {
			key := rec.BookId + `:` + strconv.Itoa(rec.ChapterNum)
			chapTitle[key] = rec
		} else {
			titles, _ := bookTitle[rec.BookId]
			titles = append(titles, rec)
			bookTitle[rec.BookId] = titles
		}
	}
	return bookTitle, chapTitle
}

func (d *DBPTextEditReader) combineHeadings(bookTitle map[string][]db.Script,
	chapTitle map[string]db.Script) []db.Script {
	var results = make([]db.Script, 0, 5000)
	var database = d.bibleId + `_DBPTEXT.db`
	var conn = db.NewDBAdapter(d.ctx, database)
	var lastBookId = ``
	var lastChapter = -1
	var scriptNum = 0
	for _, rec := range conn.SelectScripts() {
		if rec.BookId != lastBookId || rec.ChapterNum != lastChapter {
			scriptNum = 0
		}
		if rec.BookId != lastBookId {
			lastBookId = rec.BookId
			titleRec := bookTitle[rec.BookId]
			for _, title := range titleRec {
				scriptNum++
				var inp db.Script
				inp.BookId = rec.BookId
				inp.ChapterNum = rec.ChapterNum
				inp.UsfmStyle = title.UsfmStyle
				inp.ScriptNum = strconv.Itoa(scriptNum)
				inp.VerseNum = 0
				inp.VerseStr = ``
				inp.ScriptTexts = []string{title.ScriptText}
				results = append(results, inp)
			}
			lastChapter = -1
		}
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			key := rec.BookId + `:` + strconv.Itoa(rec.ChapterNum)
			head := chapTitle[key]
			scriptNum++
			var inp db.Script
			inp.BookId = rec.BookId
			inp.ChapterNum = rec.ChapterNum
			inp.UsfmStyle = head.UsfmStyle
			inp.ScriptNum = strconv.Itoa(scriptNum)
			inp.VerseNum = 0
			inp.VerseStr = ``
			inp.ScriptTexts = []string{head.ScriptText}
			results = append(results, inp)
		}
		scriptNum++
		var inp db.Script
		inp.BookId = rec.BookId
		inp.ChapterNum = rec.ChapterNum
		inp.ScriptNum = strconv.Itoa(scriptNum)
		inp.VerseNum = rec.VerseNum
		inp.VerseStr = rec.VerseStr
		inp.ScriptTexts = []string{rec.ScriptText}
		results = append(results, inp)
	}
	return results
}
