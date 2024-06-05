package read

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	"dataset/request"
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
	ctx       context.Context
	bibleId   string
	conn      db.DBAdapter
	req       request.Request // This is not needed
	testament request.Testament
}

func NewDBPTextEditReader(conn db.DBAdapter, req request.Request) DBPTextEditReader {
	var d DBPTextEditReader
	d.ctx = conn.Ctx
	d.bibleId = req.BibleId
	d.conn = conn
	d.req = req
	d.testament = req.Testament
	return d
}

func (d *DBPTextEditReader) Process() dataset.Status {
	var status dataset.Status
	testament := d.req.Testament
	var usxDB db.DBAdapter
	usxDB, status = d.createUSXEDITText(testament)
	if status.IsErr {
		return status
	}
	titleMap, chapMap, status := d.readUSXHeadings(usxDB) //testament)
	if status.IsErr {
		return status
	}
	var textDB db.DBAdapter
	textDB, status = d.createDBPText(testament)
	if status.IsErr {
		return status
	}
	records, status := d.combineHeadings(textDB, titleMap, chapMap)
	if !status.IsErr {
		d.conn.InsertScripts(records)
	}
	return status
}

func (d *DBPTextEditReader) createDBPText(testament request.Testament) (db.DBAdapter, dataset.Status) {
	var database db.DBAdapter
	var status dataset.Status
	var otMediaId string
	var ntMediaId string
	if testament.OT || len(testament.OTBooks) > 0 {
		otMediaId = d.bibleId + `O_ET`
	}
	if testament.NT || len(testament.NTBooks) > 0 {
		ntMediaId = d.bibleId + `N_ET`
	}
	files, status := input.DBPDirectory(d.ctx, d.bibleId, request.TextPlainEdit, otMediaId,
		ntMediaId, testament)
	if status.IsErr {
		return database, status
	}
	database = db.NewDBAdapter(d.ctx, ":memory:")
	textAdapter := NewDBPTextReader(database, d.req.Testament)
	status = textAdapter.ProcessFiles(files)
	return database, status
}

func (d *DBPTextEditReader) createUSXEDITText(testament request.Testament) (db.DBAdapter, dataset.Status) {
	var database db.DBAdapter
	var status dataset.Status
	files, status := input.DBPDirectory(d.ctx, d.bibleId, request.TextUSXEdit, d.bibleId+`O_ET-usx`,
		d.bibleId+`N_ET-usx`, testament)
	if status.IsErr {
		return database, status
	}
	database = db.NewDBAdapter(d.ctx, ":memory:")
	usx := NewUSXParser(database)
	status = usx.ProcessFiles(files)
	return database, status
}

func (d *DBPTextEditReader) readUSXHeadings(conn4 db.DBAdapter) (map[string][]db.Script, map[string]db.Script, dataset.Status) {
	var bookTitle = make(map[string][]db.Script)
	var chapTitle = make(map[string]db.Script)
	var status dataset.Status
	var records []db.Script
	records, status = conn4.SelectScriptHeadings()
	conn4.Close()
	if !status.IsErr {
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
	}
	return bookTitle, chapTitle, status
}

func (d *DBPTextEditReader) combineHeadings(conn db.DBAdapter, bookTitle map[string][]db.Script,
	chapTitle map[string]db.Script) ([]db.Script, dataset.Status) {
	var results = make([]db.Script, 0, 5000)
	var lastBookId = ``
	var lastChapter = -1
	var scriptNum = 0
	var records, status = conn.SelectScripts()
	conn.Close()
	if status.IsErr {
		return results, status
	}
	for _, rec := range records {
		if d.testament.HasOT(rec.BookId) || d.testament.HasNT(rec.BookId) {
			if rec.BookId != lastBookId || rec.ChapterNum != lastChapter {
				scriptNum = 1
				var inp db.Script
				inp.BookId = rec.BookId
				inp.ChapterNum = rec.ChapterNum
				inp.ScriptNum = strconv.Itoa(scriptNum)
				inp.VerseNum = 0
				inp.VerseStr = `0`
				if rec.BookId != lastBookId {
					lastBookId = rec.BookId
					titleRec := bookTitle[rec.BookId]
					for _, title := range titleRec {
						inp.UsfmStyle = title.UsfmStyle
						if len(inp.ScriptTexts) > 0 {
							inp.ScriptTexts = append(inp.ScriptTexts, ` `)
						}
						inp.ScriptTexts = append(inp.ScriptTexts, title.ScriptText)
					}
				}
				lastChapter = rec.ChapterNum
				key := rec.BookId + `:` + strconv.Itoa(rec.ChapterNum)
				head := chapTitle[key]
				scriptNum++
				if inp.UsfmStyle == `` {
					inp.UsfmStyle = head.UsfmStyle
				}
				if len(inp.ScriptTexts) > 0 {
					inp.ScriptTexts = append(inp.ScriptTexts, ` `)
				}
				inp.ScriptTexts = append(inp.ScriptTexts, head.ScriptText)
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
	}
	return results, status
}
