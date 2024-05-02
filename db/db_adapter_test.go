package db

import (
	"context"
	"testing"
)

func TestNewerDBAdapter(t *testing.T) {
	ctx := context.Background()
	conn := NewerDBAdapter(ctx, true, `GaryG`, `TestNewerDBAdapter`)
	count, status := conn.CountScriptRows()
	if status.IsErr {
		t.Fatal(status)
	}
	if count != 0 {
		t.Fatal(`Tables should be zero rows`, count)
	}
}
