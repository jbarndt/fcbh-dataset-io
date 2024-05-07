package read

import (
	"context"
	"dataset/db"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestScriptReader(t *testing.T) {
	bibleId := `ATIWBT`
	database := bibleId + "_SCRIPT.db"
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(context.Background(), database)
	script := NewScriptReader(conn)
	filename := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, bibleId+`N2ST.xlsx`)
	fmt.Println(`Filename:`, filename)
	status := script.Read(filename)
	if status.IsErr {
		t.Fatal(status)
	}
	//}
	conn.Close()
}
