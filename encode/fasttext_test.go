package encode

import (
	"context"
	"dataset/db"
	"testing"
)

func TestFastText(t *testing.T) {
	var ctx = context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	var fast = NewFastText(ctx, conn)
	fast.Process()
}
