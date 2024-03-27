package read

import (
	"dataset_io"
	"dataset_io/db"
	"testing"
)

func TestDBPEditTextReader(t *testing.T) {
	var bibleId = `ATIWBT`
	var database = bibleId + `_EDITTEXT.db`
	db.DestroyDatabase(database)
	var db1 = db.NewDBAdapter(database)
	reader := NewDBPTextEditReader(bibleId, db1)
	reader.Process(dataset_io.NT)
}
