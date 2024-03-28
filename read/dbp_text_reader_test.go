package read

import (
	"dataset"
	"dataset/db"
	"testing"
)

func TestDBPTextReader(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_DBPTEXT.db`
	db.DestroyDatabase(database)
	var db1 = db.NewDBAdapter(database)
	textAdapter := NewDBPTextReader(db1)
	textAdapter.ProcessDirectory(bibleId, dataset.NT)
}
