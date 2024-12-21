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
	var database = "ENGWEB_align_mp3.db"
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", database)
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn, "eng", "")
	audioDir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), "ENGWEB", "ENGWEBN2DA-mp3-64")
	faVerses, filenameMap, status := calc.Process(audioDir)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(filenameMap)
	calc.countErrors(faVerses)
	fmt.Println(len(faVerses))
}
