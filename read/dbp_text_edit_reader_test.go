package read

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestDBPEditTextReader(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_EDITTEXT.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	var db1 = db.NewDBAdapter(ctx, database)
	reader := NewDBPTextEditReader(bibleId, db1)
	reader.Process(request.Testament{NT: true})
}
