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
	//file.BookId = "PHM"
	files = append(files, file)
	status = fa.ProcessFiles(files)
	if status.IsErr {
		t.Fatal(status)
	}
}

func TestMMSFA_processPyOutput(t *testing.T) {
	ctx := context.Background()
	user, _ := fetch.GetTestUser()
	conn, status := db.NewerDBAdapter(ctx, false, user.Username, "PlainTextEditScript_ENGWEB")
	if status.IsErr {
		t.Fatal(status)
	}
	fa := NewMMSFA(ctx, conn, "eng", "")
	var file input.InputFile
	file.BookId = "MRK"
	file.Chapter = 1
	file.MediaId = "ENGWEBN2DA"
	file.Directory = os.Getenv("FCBH_DATASET_FILES") + "/ENGWEB/ENGWEBN2DA-mp3-64/"
	file.Filename = "B02___01_Mark________ENGWEBN2DA.mp3"
	verses, status := conn.SelectScriptsByChapter(file.BookId, file.Chapter)
	if status.IsErr {
		t.Fatal(status)
	}
	var wordList []Word
	_, wordList, status = fa.prepareText("eng", verses)
	if status.IsErr {
		t.Fatal(status)
	}
	bytes, err := os.ReadFile("engweb_fa_out.json")
	if err != nil {
		t.Fatal(err)
	}
	fa.processPyOutput(file, wordList, string(bytes))
	scriptRows, status := conn.SelectScalarInt("select count(*) from scripts where script_begin_ts != 0.0")
	if status.IsErr {
		t.Fatal(status)
	}
	if scriptRows != 46 {
		t.Error("scriptRows is", scriptRows, "it should be 46")
	}
	wordRows, status := conn.SelectScalarInt("select count(*) from words")
	if status.IsErr {
		t.Fatal(status)
	}
	if wordRows != 883 {
		t.Error("wordRows is", wordRows, "it should be 883")
	}
}
