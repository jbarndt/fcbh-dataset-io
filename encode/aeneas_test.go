package encode

import (
	"context"
	"dataset/db"
	"dataset/read"
	"dataset/request"
	"testing"
)

func TestAeneasLines(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var language = `eng`
	var testament = request.Testament{NTBooks: []string{`MRK`}}
	testament.BuildBookMaps()
	var detail = request.Detail{Lines: true}
	files, status := read.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	aeneas := NewAeneas(ctx, conn, bibleId, language, detail)
	status = aeneas.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
}

func TestAeneasWords(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var language = `eng`
	var testament = request.Testament{NT: true}
	testament.BuildBookMaps()
	var detail = request.Detail{Words: true}
	files, status := read.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	aeneas := NewAeneas(ctx, conn, bibleId, language, detail)
	status = aeneas.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
}
