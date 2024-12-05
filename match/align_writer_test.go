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
	faVerses, filenameMap, status := calc.Process()
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(len(faVerses), len(filenameMap))
	writer := NewAlignWriter(ctx)
	filename, status := writer.WriteReport(dataset, faVerses, filenameMap)
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}
