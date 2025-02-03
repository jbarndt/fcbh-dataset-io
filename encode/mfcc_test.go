package encode

import (
	"context"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/db"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/decode_yaml/request"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/input"
	"testing"
)

func TestMFCCLines(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var testament = request.Testament{NTBooks: []string{`MRK`}}
	testament.BuildBookMaps()
	var detail = request.Detail{Lines: true}
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status != nil {
		t.Error(status)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	conn.DeleteMFCCs()
	mfcc := NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status != nil {
		t.Error(status)
	}
	count, _ := conn.CountScriptMFCCRows()
	if count != 678 {
		t.Error(`Script count should be 678`, count)
	}
	count, _ = conn.CountWordMFCCRows()
	if count != 0 {
		t.Error(`Word count should be 0`, count)
	}
}

func TestMFCCWords(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var testament = request.Testament{NTBooks: []string{`MRK`}}
	testament.BuildBookMaps()
	var detail = request.Detail{Words: true}
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status != nil {
		t.Error(status)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	mfcc := NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status != nil {
		t.Error(status)
	}
	count, _ := conn.CountScriptMFCCRows()
	if count != 678 {
		t.Error(`Script count should be 1`, count)
	}
	count, _ = conn.CountWordMFCCRows()
	if count != 0 {
		t.Error(`Word count should be 0`, count)
	}
}
