package db

import (
	"context"
	"testing"
)

func TestNewerDBAdapter(t *testing.T) {
	ctx := context.Background()
	conn, status1 := NewerDBAdapter(ctx, true, `GaryG`, `TestNewerDBAdapter`)
	if status1 != nil {
		t.Fatal(status1)
	}
	count, status := conn.CountScriptRows()
	if status != nil {
		t.Fatal(status)
	}
	if count != 0 {
		t.Fatal(`Tables should be zero rows`, count)
	}
}
