package encode

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/fetch"
	"dataset/input"
	log "dataset/logger"
	"dataset/read"
	"dataset/request"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 1. Load the script into a memory database
// 2. Parsing the script, find those script lines that begin a verse.
// 3. Break the script up into parts that contain 1 or more whole verse
// 4. Load the plain_text_edit script
// 5. Iterate over the script groups
// 6. Access the script timestamps each group
// 7. Set timestamps in grouped script
// 8. For each group get the corresponding verses from plain_text_edit
// 9. Process verses through aeneas
// 10. Compute the final timestamps of all verses
// 11. Update plain_text_edit timestamps

type AeneasExperiment struct {
	ctx context.Context
	ts  input.TSBucket
}

func NewAeneasExperiment(ctx context.Context) AeneasExperiment {
	var a AeneasExperiment
	a.ctx = ctx
	a.ts = input.NewTSBucket(a.ctx)
	return a
}

func (a *AeneasExperiment) Process() {
	audioMediaId := `APFCMUN2DA`
	script := a.LoadScript(audioMediaId)
	fmt.Println(`script lines`, len(script))
	textConn := a.LoadPlainTextEdit(audioMediaId)
	a.ProcessByChapter(audioMediaId, script, textConn)
}

func (a *AeneasExperiment) LoadScript(mediaId string) []db.Script {
	var results []db.Script
	var status dataset.Status
	//ts := input.NewTSBucket(a.ctx)
	key := a.ts.GetKey(input.Script, mediaId, ``, 0)
	fmt.Println(`key:`, key)
	filePath := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), mediaId[:6], mediaId[:8]+`ST.xlsx`)
	a.ts.DownloadObject(input.TSBucketName, key, filePath)
	fmt.Println(`path`, filePath)
	var conn = db.NewDBAdapter(a.ctx, ":memory:")
	reader := read.NewScriptReader(conn)
	reader.Read(filePath)
	results, status = conn.SelectScripts()
	if status.IsErr {
		panic(status)
	}
	conn.Close()
	fmt.Println(`count:`, len(results))
	return results
}

func (a *AeneasExperiment) LoadPlainTextEdit(audioMediaId string) db.DBAdapter {
	user, _ := fetch.GetTestUser()
	bibleId := audioMediaId[:6]
	conn, status := db.NewerDBAdapter(a.ctx, true, user.Username, `Plain_Text_Edit_`+bibleId)
	if status.IsErr {
		panic(status)
	}
	var req request.Request
	req.BibleId = bibleId
	req.Testament = request.Testament{NT: true}
	reader := read.NewDBPTextEditReader(conn, req)
	reader.Process()
	return conn
}

func (a *AeneasExperiment) ProcessByChapter(audioId string, scripts []db.Script, conn db.DBAdapter) {
	var results []db.Script
	var lastBookId = ``
	var lastChapter = -1
	for _, scp := range scripts {
		if scp.BookId != lastBookId || scp.ChapterNum != lastChapter {
			if len(results) > 0 {
				grouped := a.GroupScriptByVerse(results)
				grouped = a.GetTimestamps(audioId, grouped)
				a.GetPlainTextEdit(conn, grouped)
			}
			results = nil
			lastBookId = scp.BookId
			lastChapter = scp.ChapterNum
		}
		results = append(results, scp)
	}
}

func (a *AeneasExperiment) GroupScriptByVerse(scripts []db.Script) []db.Script {
	var results = make([]db.Script, 0, len(scripts))
	var rec db.Script
	for _, scp := range scripts {
		trimmed := strings.TrimSpace(scp.ScriptText)
		if trimmed[0] == '{' || scp.VerseStr == `0` {
			if rec.BookId != `` {
				results = append(results, rec)
			}
			rec = scp
		}
		rec.ScriptNum = scp.ScriptNum
		rec.VerseEnd = scp.VerseStr
		rec.ScriptTexts = append(rec.ScriptTexts, scp.ScriptText)
	}
	if rec.BookId != `` {
		results = append(results, rec)
	}
	for i, scp := range results {
		results[i].ScriptText = strings.Join(scp.ScriptTexts, "")
	}
	return results
}

func (a *AeneasExperiment) GetTimestamps(audioId string, scripts []db.Script) []db.Script {
	timestamps := a.ts.GetTimestamps(input.ScriptTS, audioId, scripts[0].BookId, scripts[0].ChapterNum)
	fmt.Println(`timestamps:`, timestamps)
	var timeMap = make(map[string]db.Timestamp)
	for _, t := range timestamps {
		timeMap[t.VerseStr] = t
	}
	for i, scp := range scripts {
		ts, ok := timeMap[scp.VerseStr]
		if ok {
			scripts[i].ScriptBeginTS = ts.BeginTS
		} else {
			log.Warn(a.ctx, `No timestamp found for verse `, scp.BookId, scp.ChapterNum, scp.VerseStr)
		}
		ts, ok = timeMap[scp.VerseEnd]
		if ok {
			scripts[i].ScriptEndTS = ts.EndTS
		} else {
			log.Warn(a.ctx, `No timestamp found for verse `, scp.BookId, scp.ChapterNum, scp.VerseEnd)
		}
	}
	return scripts
}

func (a *AeneasExperiment) GetPlainTextEdit(conn db.DBAdapter, groups []db.Script) {
	scripts, status := conn.SelectScriptsByChapter(groups[0].BookId, groups[0].ChapterNum)
	if status.IsErr {
		panic(status)
	}
	var beginMap = make(map[string]float64)
	var endMap = make(map[string]float64)
	for _, grp := range groups {
		beginMap[grp.VerseStr] = grp.ScriptBeginTS
		endMap[grp.VerseEnd] = grp.ScriptEndTS
	}
	for i, scp := range scripts {
		scripts[i].ScriptBeginTS = -1.0
		scripts[i].ScriptEndTS = -1.0
		ts, ok := beginMap[scp.VerseStr]
		if ok {
			scripts[i].ScriptBeginTS = ts
		}
		ts, ok = endMap[scp.VerseStr]
		if ok {
			scripts[i].ScriptEndTS = ts
		}
	}
	for _, part := range groups {
		if part.VerseStr != part.VerseEnd {
			var recs []db.Script
			for _, scp := range scripts {
				if scp.VerseStr >= part.VerseStr && scp.VerseEnd <= part.VerseEnd {
					recs = append(recs, scp)
				}
			}
			a.ProcessAeneas(recs)
		}
	}
}

func (a *AeneasExperiment) ProcessAeneas(recs []db.Script) {
	for _, scp := range recs {
		fmt.Println(`Process in Aeneas`, scp)
	}
}

func (a *AeneasExperiment) CountChars(scripts []db.Script) int {
	var results int
	for _, scp := range scripts {
		results += len(scp.ScriptText)
	}
	return results
}

func (a *AeneasExperiment) DumpScript(scripts []db.Script, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, scp := range scripts {
		//_, err = fmt.Fprintln(file, scp.BookId, scp.ChapterNum, scp.VerseStr, scp.ScriptNum, scp.ScriptText)
		_, err = fmt.Fprintln(file, scp.ScriptText)
		if err != nil {
			panic(err)
		}
	}
}
