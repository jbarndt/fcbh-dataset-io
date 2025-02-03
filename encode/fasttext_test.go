package encode

import (
	"context"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/db"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/decode_yaml/request"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/input"
	"github.com/faithcomesbyhearing/fcbh-dataset-io/read"
	"testing"
)

func TestFastText(t *testing.T) {
	var ctx = context.Background()
	db.DestroyDatabase(`ENGWEB_DBPTEXT.db`)
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	bibleId := `ENGWEB`
	testament := request.Testament{NT: true}
	files, status := input.DBPDirectory(ctx, bibleId, `text_plain`, `ENGWEBO_ET`,
		`ENGWEBN_ET`, testament)
	reader := read.NewDBPTextReader(conn, testament)
	status = reader.ProcessFiles(files)
	if status != nil {
		t.Error(status)
	}
	words := read.NewWordParser(conn)
	status = words.Parse()
	if status != nil {
		t.Error(status)
	}
	var fast = NewFastText(ctx, conn)
	status = fast.Process()
	if status != nil {
		t.Error(status)
	}
}
