package read

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestDBPEditTextReader(t *testing.T) {
	var req request.Request
	req.BibleId = `ENGWEB`
	req.Testament = request.Testament{OTBooks: []string{`GEN`, `EXO`}, NTBooks: []string{`MAT`, `MRK`}}
	req.Testament.BuildBookMaps()
	var database = req.BibleId + `_EDITTEXT.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	var db1 = db.NewDBAdapter(ctx, database)
	reader := NewDBPTextEditReader(db1, req)
	status := reader.Process()
	if status.IsErr {
		t.Error(status.Message)
	}
	count, status := db1.CountScriptRows()
	if status.IsErr {
		t.Error(status.Message)
	}
	if count != 4629 {
		t.Error(`Expected count to be 4629`, count)
	}
}
