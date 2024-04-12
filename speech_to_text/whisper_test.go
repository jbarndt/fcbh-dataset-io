package speech_to_text

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestWhisper(t *testing.T) {
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-opus16`
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, database)
	var whisp = NewWhisper(bibleId, conn, `small`)
	whisp.ProcessDirectory(filesetId, request.Testament{NT: true})
	conn.Close()
}
