package match

import (
	"context"
	"dataset/db"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewAlignErrorCalc(t *testing.T) {
	ctx := context.Background()
	//var database = "N2YPM_JMD.db"
	var database = "ENGWEB.db"
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", database)
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn)
	faVerses, filenameMap, status := calc.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(filenameMap)
	calc.countErrors(faVerses)
	fmt.Println(len(faVerses))
}

// TestComputeAvgIntervals is not a Test, but a routine to compute mean and std dev of time between timestamps

/*
This was useful it updated the database with a silence value
func TestUpdateCharSilence(t *testing.T) {
	ctx := context.Background()
	var database = "ENGWEB.db"
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", database)
	conn := db.NewDBAdapter(ctx, dbPath)
	defer conn.Close()
	chars, _ := conn.SelectFACharTimestamps()
	query := `UPDATE chars SET silence = ? WHERE char_id = ?`
	tx, err := conn.Begin()
	if err != nil {
		t.Fatal(err)
	}
	stmt, err := tx.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < len(chars)-1; i++ {
		silence := chars[i+1].charBeginTS - chars[i].charEndTS
		_, err = stmt.Exec(silence, chars[i].charId)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

*/
