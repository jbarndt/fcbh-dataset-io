package speech_to_text

import (
	"context"
	"dataset/db"
	"dataset/read"
	"dataset/request"
	"testing"
)

func TestWhisper(t *testing.T) {
	ctx := context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-opus16`
	testament := request.Testament{NTBooks: []string{`TIT`, `PHM`, `3JN`}}
	testament.BuildBookMaps()
	files, status := read.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(ctx, database)
	var whisp = NewWhisper(bibleId, conn, `tiny`)
	status = whisp.ProcessFiles(files)

	if status.IsErr {
		t.Error(status.Message)
	}
	count, status := conn.CountScriptRows()
	if count != 1 {
		t.Error(`CountScriptRows count != 1`, count)
	}
	conn.Close()
}
