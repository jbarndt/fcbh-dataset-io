package speech_to_text

import (
	"context"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/db"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/decode_yaml/request"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/input"
	"testing"
)

func TestWhisper(t *testing.T) {
	ctx := context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-mp3-64`
	testament := request.Testament{NTBooks: []string{`TIT`, `PHM`, `3JN`}}
	testament.BuildBookMaps()
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status != nil {
		t.Error(status)
	}
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(ctx, database)
	var whisp = NewWhisper(bibleId, conn, `tiny`, `en`)
	status = whisp.ProcessFiles(files)

	if status != nil {
		t.Error(status)
	}
	count, status := conn.CountScriptRows()
	if count != 120 {
		t.Error(`CountScriptRows count != 120`, count)
	}
	conn.Close()
}
