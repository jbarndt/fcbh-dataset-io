package speech_to_text

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/fetch"
	"dataset/input"
	"dataset/read"
	"dataset/request"
	"fmt"
	"testing"
)

func TestWhisperVs(t *testing.T) {
	ctx := context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-mp3-64`
	testament := request.Testament{NTBooks: []string{`TIT`, `PHM`, `3JN`}}
	//testament := request.Testament{NTBooks: []string{`3JN`}}
	testament.BuildBookMaps()
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Fatal(status)
	}
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(ctx, database)
	loadPlainText(bibleId, conn, testament, t)
	loadTimestamps(filesetId, conn, testament, t)
	newConn, status := conn.CopyDatabase(`_STT`)
	if status.IsErr {
		t.Fatal(status)
	}
	var whisp = NewWhisperVs(bibleId, newConn, `tiny`)
	status = whisp.ProcessFiles(files)
	if status.IsErr {
		t.Fatal(status)
	}
	count, status := newConn.CountScriptRows()
	if count != 90 {
		t.Error(`CountScriptRows count != 90`, count)
	}
	newConn.Close()
}

func loadPlainText(bibleId string, conn db.DBAdapter,
	testament request.Testament, t *testing.T) {
	var status dataset.Status
	req := request.Request{}
	req.BibleId = bibleId
	req.Testament = testament
	parser := read.NewDBPTextEditReader(conn, req)
	status = parser.Process()
	if status.IsErr {
		t.Error(status.Message)
	}
}

func loadTimestamps(filesetId string, conn db.DBAdapter,
	testament request.Testament, t *testing.T) {
	api := fetch.NewAPIDBPTimestamps(conn, filesetId)
	ok, status := api.LoadTimestamps(testament)
	if status.IsErr {
		t.Error(status)
	}
	fmt.Println("Timestamps OK ", ok)
}
