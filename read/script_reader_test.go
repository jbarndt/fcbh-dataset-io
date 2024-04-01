package read

import (
	"context"
	"dataset/db"
	"testing"
)

func TestScriptReader(t *testing.T) {
	bibleId := `ATIWBT`
	database := bibleId + "_SCRIPT.db"
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(context.Background(), database)
	script := NewScriptReader(conn)
	filename := script.FindFile(bibleId)
	script.Read(filename)
	conn.Close()
}
