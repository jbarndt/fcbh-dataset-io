package speech_to_text

import (
	"context"
	"dataset"
	"dataset/db"
	"testing"
)

func TestWhisper(t *testing.T) {
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-opus16`
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	ctx := context.Background()
	conn := db.NewDBAdapter(ctx, database)
	var whisp = NewWhisper(bibleId, conn)
	whisp.ProcessDirectory(filesetId, dataset.NT)
	conn.Close()
}
