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
	for i, row := range scripts {
		fmt.Println(i, row)
		if i > 5 {
			break
		}
	}
	fmt.Println(len(scripts))
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
