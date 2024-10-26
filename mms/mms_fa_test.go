package mms

import (
	"context"
	"dataset/db"
	"dataset/fetch"
	"dataset/input"
	"os"
	"testing"
)

func TestMMSFA_ProcessFiles(t *testing.T) {
	ctx := context.Background()
	user, _ := fetch.GetTestUser()
	conn, status := db.NewerDBAdapter(ctx, false, user.Username, "PlainTextEditScript_ENGWEB")
	fa := NewMMSFA(ctx, conn, "eng", "")
	var files []input.InputFile
	var file input.InputFile
	file.BookId = "MRK"
	file.Chapter = 1
	file.MediaId = "ENGWEBN2DA"
	file.Directory = os.Getenv("FCBH_DATASET_FILES") + "/ENGWEB/ENGWEBN2DA-mp3-64/"
	file.Filename = "B02___01_Mark________ENGWEBN2DA.mp3"
	//file.MediaId = "ENGESVN1DA"
	//file.Directory = os.Getenv("FCBH_DATASET_FILES") + "/ENGESV/ENGESVN1DA/"
	//file.Filename = "B02___01_Mark________ENGESVN1DA.mp3"
	files = append(files, file)
	status = fa.ProcessFiles(files)
	if status.IsErr {
		t.Fatal(status)
	}
}
