package match

import (
	"context"
	"dataset/db"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAlignWriter(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", "N2YPM_JMD.db")
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn)
	verses, status := calc.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	writer := NewAlignWriter(ctx)
	filename, status := writer.WriteReport("N2YPM_JMD", verses)
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", "N2YPM_JMD.html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}
