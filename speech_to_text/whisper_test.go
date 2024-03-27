package speech_to_text

import (
	"dataset_io"
	"dataset_io/db"
	"testing"
)

func TestWhisper(t *testing.T) {
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-opus16`
	var database = bibleId + `_WHISPER.db`
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(database)
	var whisp = NewWhisper(bibleId, conn)
	whisp.ProcessDirectory(filesetId, dataset_io.NT)
	conn.Close()
}
