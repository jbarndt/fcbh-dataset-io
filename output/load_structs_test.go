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

func TestLoadScriptSrtuct(t *testing.T) {
	var ctx = context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	scripts := LoadScriptStruct(conn)
	var script Script
	metaScript := ReflectStruct(script)
	newMetaScript := FindActiveScriptCols(scripts, metaScript)
	numMFCC := FindNumMFCC(scripts)
	for i, meta := range newMetaScript {
		if meta.Name == `MFCC` {
			newMetaScript[i].Cols = numMFCC
		}
	}
	scripts2 := NormalizeMFCC(scripts, numMFCC)
	scripts3 := PadRows(scripts2, numMFCC)
	fmt.Println(newMetaScript)
	fmt.Println("length scripts", len(scripts))
	fmt.Println("length normalized", len(scripts))
	fmt.Println("length padded", len(scripts3))
	filename := WriteScriptCSV(scripts3, newMetaScript)
	fmt.Println(`CSV Script`, filename)
	filename2 := WriteScriptJSON(scripts3, newMetaScript)
	fmt.Println(`JSON Script`, filename2)
}

func prepareTimestampAndFMCCData(conn db.DBAdapter, bibleId string, filesetId string, t *testing.T) {
	ctx := context.Background()
	api := fetch.NewAPIDBPTimestamps(conn, filesetId)
	testament := request.Testament{NTBooks: []string{`MAT`, `MRK`}}
	testament.BuildBookMaps()
	_, status := api.LoadTimestamps(testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	var detail = request.Detail{Lines: true}
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	mfcc := encode.NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
}

//func TestReflectStruct(t *testing.T) {
//	ReflectScriptStruct()
//}
