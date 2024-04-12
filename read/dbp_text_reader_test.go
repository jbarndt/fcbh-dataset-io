package read

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestDBPTextReader(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_DBPTEXT.db`
	db.DestroyDatabase(database)
	var db1 = db.NewDBAdapter(context.Background(), database)
	textAdapter := NewDBPTextReader(db1)
	textAdapter.ProcessDirectory(bibleId, request.Testament{NT: true})
}
