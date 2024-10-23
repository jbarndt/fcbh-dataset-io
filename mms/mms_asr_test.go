package mms

import (
	"context"
	"dataset/db"
	"dataset/input"
	"os"
	"testing"
)

func TestMMSASR_ProcessFiles(t *testing.T) {
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, ":memory:")
	asr := NewMMSASR(ctx, conn, "eng", "")
	var files []input.InputFile
	var file input.InputFile
	file.MediaId = "ENGWEBN2DA"
	file.BookId = "MRK"
	file.Chapter = 1
	file.Directory = os.Getenv("FCBH_DATASET_FILES") + "/ENGWEB/ENGWEBN2DA-mp3-64/"
	file.Filename = "B02___01_Mark________ENGWEBN2DA.mp3"
	files = append(files, file)
	status := asr.ProcessFiles(files)
	if status.IsErr {
		t.Fatal(status)
	}
}
