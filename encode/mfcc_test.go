package encode

import (
	"context"
	"dataset/db"
	"dataset/input"
	"dataset/request"
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
	if status.IsErr {
		t.Error(status.Message)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	conn.DeleteMFCCs()
	mfcc := NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
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
	if status.IsErr {
		t.Error(status.Message)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	mfcc := NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
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
