package encode

import (
	"context"
	"dataset"
	"dataset/db"
	"testing"
)

func TestMFCCLines(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//files, status := ReadDirectory(ctx, bibleId, filesetId)
	mfcc := NewMFCC(ctx, conn, bibleId, filesetId)
	mfcc.Process(dataset.LINES)
}
