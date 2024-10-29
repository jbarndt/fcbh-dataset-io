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

func TestMMSFA_processPyOutput(t *testing.T) {
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, "TestMMSFA_processPyOutput")
	fa := NewMMSFA(ctx, conn, "eng", "")
	//var files []input.InputFile
	var file input.InputFile
	file.BookId = "MRK"
	file.Chapter = 1
	file.MediaId = "ENGWEBN2DA"
	file.Directory = os.Getenv("FCBH_DATASET_FILES") + "/ENGWEB/ENGWEBN2DA-mp3-64/"
	file.Filename = "B02___01_Mark________ENGWEBN2DA.mp3"
	bytes, err := os.ReadFile("engweb_fa_out.json")
	if err != nil {
		t.Fatal(err)
	}
	fa.processPyOutput(file, string(bytes))
	scriptRows, status := conn.CountScriptRows()
	if status.IsErr {
		t.Fatal(status)
	}
	if scriptRows != 46 {
		t.Error("scriptRows is", scriptRows, "it should be 46")
	}
	wordRows, status := conn.CountWordRows()
	if status.IsErr {
		t.Fatal(status)
	}
	if wordRows != 884 {
		t.Error("wordRows is", wordRows, "it should be 884")
	}
}
