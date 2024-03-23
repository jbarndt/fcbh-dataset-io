package read

import (
	"dataset_io/db"
	"testing"
)

func TestScriptReader(t *testing.T) {
	bibleId := `ATIWBT`
	database := bibleId + "_SCRIPT.db"
	db.DestroyDatabase(database)
	db := db.NewDBAdapter(database)
	script := NewScriptReader(db)
	filename := script.FindFile(bibleId)
	script.Read(filename)
	db.Close()
}
