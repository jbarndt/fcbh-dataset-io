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
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", "N2YPM_JMD.db")
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn)
	verses, status := calc.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	for _, vs := range verses {
		fmt.Printf("%v\n", vs)
	}
	fmt.Println(len(verses))
}
