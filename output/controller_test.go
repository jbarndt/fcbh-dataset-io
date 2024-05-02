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

func TestPrepareScripts(t *testing.T) {
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	var out = NewOutput(ctx, conn, true, true)
	structs, meta := out.PrepareScripts()
	fmt.Println("Loaded Scripts", len(structs))
	filename, status := out.WriteCSV(structs, meta)
	if status.IsErr {
		t.Error(status)
	}
	fmt.Println("CoSV File", filename)
	filename, status = out.WriteJSON(structs, meta)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("JSON File", filename)
}

func TestPrepareWords(t *testing.T) {
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	var out = NewOutput(ctx, conn, true, true)
	structs, meta := out.PrepareWords()
	fmt.Println("Loaded Scripts", len(structs))
	filename, status := out.WriteCSV(structs, meta)
	if status.IsErr {
		t.Error(status)
	}
	fmt.Println("CSV File", filename)
	filename, status = out.WriteJSON(structs, meta)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println("JSON File", filename)
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
