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
	dataset := "N2ENGWEB"
	dbDir := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match")
	conn := db.NewDBAdapter(ctx, filepath.Join(dbDir, "N2ENGWEB.db"))
	asrConn := db.NewDBAdapter(ctx, filepath.Join(dbDir, "N2ENGWEB_audio.db"))
	calc := NewAlignErrorCalc(ctx, conn, asrConn, "eng", "")
	audioDir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), "ENGWEB", "ENGWEBN2DA-mp3-64")
	faLines, filenameMap, status := calc.Process(audioDir)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(len(faLines), len(filenameMap))
	writer := NewAlignWriter(ctx)
	filename, status := writer.WriteReport(dataset, faLines, filenameMap)
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}
