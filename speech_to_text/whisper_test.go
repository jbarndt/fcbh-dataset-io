package speech_to_text

import (
	"dataset_io"
	"dataset_io/db"
	"testing"
)

func TestWhisper(t *testing.T) {
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA-opus16`
	var conn = db.NewDBAdapter(bibleId + `_WHISPER.db`)
	var whisp = NewWhisper(conn)
	whisp.ProcessDirectory(bibleId, filesetId, dataset_io.NT)
	conn.Close()
}
