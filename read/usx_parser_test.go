package read

import (
	"dataset_io"
	"dataset_io/db"
	"testing"
)

func TestUSXParser(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_USXEDIT.db`
	db.DestroyDatabase(database)
	var conn = db.NewDBAdapter(database)
	ReadUSXEdit(conn, bibleId, dataset.NT)
}
