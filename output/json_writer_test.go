package output

import (
	"context"
	"dataset/db"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/floats"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestJSONWriterScript(t *testing.T) {
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	structs, meta := PrepareScripts(conn, false, false)
	fmt.Println("Loaded Scripts", len(structs))
	filename := WriteJSON(structs, meta)
	fileRecs := readJSONScript(filename, t)
	fmt.Println(len(fileRecs))
	dbRecs := LoadScriptStruct(conn)
	//mfccDB := extractScriptMFCC(scripts)
	fmt.Println("Written CSV", filename)
	compareScript(dbRecs, fileRecs, t)
	fmt.Println("Written CSV", filename)
}

/*
	func TestJSONWriterWord(t *testing.T) {
		ctx := context.Background()
		var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
		//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
		structs, meta := PrepareWords(conn, false, false)
		fmt.Println("Loaded Scripts", len(structs))
		filename := WriteJSON(structs, meta)
		mfccFile := readJSONScript(filename, t)
		words := LoadWordStruct(conn)
		mfccDB := extractWordEncAndMFCC(words)
		compare(mfccDB, mfccFile, t)
		fmt.Println("Written CSV", filename)
	}
*/
func readJSONScript(filename string, t *testing.T) []Script {
	var results = make([]Script, 0, 10000)
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	var records []map[string]any
	err = json.Unmarshal(content, &records)
	if err != nil {
		t.Error(err)
	}
	var sc Script
	for _, rec := range records {
		if rec["script_id"] != nil {
			if sc.ScriptId != 0 {
				results = append(results, sc)
			}
			sc = Script{}
			sc.BookId = rec["book_id"].(string)
			sc.ChapterNum = int(rec["chapter_num"].(float64))
			sc.Reference = rec["reference"].(string)
			sc.ScriptBeginTS = rec["script_begin_ts"].(float64)
			sc.ScriptEndTS = rec["script_end_ts"].(float64)
			sc.ScriptId = int(rec["script_id"].(float64))
			sc.ScriptNum = rec["script_num"].(string)
			sc.ScriptText = rec["script_text"].(string)
			sc.VerseStr = rec["verse_str"].(string)
		}
		var mfcc = make([]float64, 300)
		var cols = 0
		for key, val := range rec {
			if strings.HasPrefix(key, "mfcc") {
				num, err := strconv.Atoi(key[4:])
				if err != nil {
					t.Error(err)
				}
				mfcc[num] = val.(float64)
				if num > cols {
					cols = num
				}
			}
		}
		sc.MFCCRows++
		sc.MFCCCols = cols + 1
		sc.MFCC = append(sc.MFCC, mfcc[:sc.MFCCCols])
	}
	results = append(results, sc)
	return results
}

func compareScript(dbRecs []Script, fileRecs []Script, t *testing.T) {
	if len(dbRecs) != len(fileRecs) {
		t.Error(`dbRecs has length`, len(dbRecs), `fileRecs has length`, len(fileRecs))
	}
	fmt.Println(`Compare `, len(dbRecs), `records`)
	var rowDiffCount = 0
	for i, db := range dbRecs {
		file := fileRecs[i]
		if db.BookId != file.BookId {
			t.Error(`BookId mismatch`, db.BookId, file.BookId)
		}
		if db.ChapterNum != file.ChapterNum {
			t.Error(`ChapterNum mismatch`, db.ChapterNum, file.ChapterNum)
		}
		if db.Reference != file.Reference {
			t.Error(`Reference mismatch`, db.Reference, file.Reference)
		}
		if db.ScriptBeginTS != file.ScriptBeginTS {
			t.Error(`ScriptBeginTS mismatch`, db.ScriptBeginTS, file.ScriptBeginTS)
		}
		if db.ScriptEndTS != file.ScriptEndTS {
			t.Error(`ScriptEndTS mismatch`, db.ScriptEndTS, file.ScriptEndTS)
		}
		if db.ScriptId != file.ScriptId {
			t.Error(`ScriptId mismatch`, db.ScriptId, file.ScriptId)
		}
		if db.ScriptNum != file.ScriptNum {
			t.Error(`ScriptNum mismatch`, db.ScriptNum, file.ScriptNum)
		}
		if db.ScriptText != file.ScriptText {
			t.Error(`ScriptText mismatch`, db.ScriptText, file.ScriptText)
		}
		if db.VerseStr != file.VerseStr {
			t.Error(`VerseStr mismatch`, db.VerseStr, file.VerseStr)
		}
		for row := 0; row < db.MFCCRows; row++ {
			if !floats.EqualApprox(db.MFCC[row], file.MFCC[row], 0.00001) {
				t.Error(`MFCC mismatch in row`, row)
			}
		}
	}
	fmt.Println("rows", len(dbRecs), `differences`, rowDiffCount)
}
