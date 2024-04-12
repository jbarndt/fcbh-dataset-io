package read

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestUSXParser(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_USXEDIT.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, database)
	ReadUSXEdit(conn, bibleId, request.Testament{NT: true})
}
