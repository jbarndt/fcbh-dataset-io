package read

import (
	"context"
	"dataset"
	"dataset/db"
	"testing"
)

func TestUSXParser(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_USXEDIT.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, database)
	ReadUSXEdit(conn, bibleId, dataset.NT)
}
