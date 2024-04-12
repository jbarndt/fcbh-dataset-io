package encode

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestMFCCLines(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	mfcc := NewMFCC(ctx, conn, bibleId, filesetId)
	mfcc.Process(request.Detail{Lines: true}, 7)
}
