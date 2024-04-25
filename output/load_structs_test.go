package output

import (
	"context"
	"dataset/db"
	"dataset/encode"
	"dataset/fetch"
	"dataset/input"
	"dataset/request"
	"fmt"
	"testing"
)

func TestScriptOutput(t *testing.T) {
	var ctx = context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	scripts := LoadScriptStruct(conn)
	var script Script
	metaScript := ReflectStruct(script)
	newMetaScript := FindActiveScriptCols(scripts, metaScript)
	numMFCC := FindNumScriptMFCC(scripts)
	for i, meta := range newMetaScript {
		if meta.Name == `MFCC` {
			newMetaScript[i].Cols = numMFCC
		}
	}
	hasMFCC2 := NormalizeMFCC(scripts, numMFCC)
	hasMFCC3 := PadRows(hasMFCC2, numMFCC)
	//hasMFCC4 := PointerToActual(hasMFCC3)
	fmt.Println(newMetaScript)
	fmt.Println("length scripts", len(scripts))
	fmt.Println("length normalized", len(hasMFCC2))
	fmt.Println("length padded", len(hasMFCC3))
	filename := WriteScriptCSV(scripts, newMetaScript)
	fmt.Println(`CSV Script`, filename)
	filename2 := WriteScriptJSON(scripts, newMetaScript)
	fmt.Println(`JSON Script`, filename2)
}

func TestWordOutput(t *testing.T) {
	var ctx = context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	words := LoadWordStruct(conn)
	var word Word
	metaWord := ReflectStruct(word)
	newMetaWord := FindActiveWordCols(words, metaWord)
	numMFCC := FindNumWordMFCC(words)
	for i, meta := range newMetaWord {
		if meta.Name == `MFCC` {
			newMetaWord[i].Cols = numMFCC
		}
	}
	fmt.Println(newMetaWord)
	//words2 := NormalizeWordMFCC(words, numMFCC)
	//scripts3 := PadRows(scripts2, numMFCC)
	fmt.Println("length words", len(words))
	//fmt.Println("length normalized", len(scripts))
	//fmt.Println("length padded", len(scripts3))
	//filename := WriteScriptCSV(scripts3, newMetaScript)
	//fmt.Println(`CSV Script`, filename)
	//filename2 := WriteScriptJSON(scripts3, newMetaScript)
	//fmt.Println(`JSON Script`, filename2)
}

func prepareTimestampAndFMCCData(conn db.DBAdapter, bibleId string, filesetId string, t *testing.T) {
	ctx := context.Background()
	api := fetch.NewAPIDBPTimestamps(conn, filesetId)
	testament := request.Testament{NTBooks: []string{`MRK`}}
	testament.BuildBookMaps()
	_, status := api.LoadTimestamps(testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	aeneas := encode.NewAeneas(ctx, conn, bibleId, `eng`, request.Detail{Words: true})
	status = aeneas.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
	var detail = request.Detail{Lines: true, Words: true}
	mfcc := encode.NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
}

//func TestReflectStruct(t *testing.T) {
//	ReflectScriptStruct()
//}
