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
	//var dataset = "N2YPM_JMD"
	var dataset = "ENGWEB"
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".db")
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn)
	verses, status := calc.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	writer := NewAlignWriter(ctx)
	filename, status := writer.WriteReport(dataset, verses)
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}
