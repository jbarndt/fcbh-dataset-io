package input

import (
	"context"
	"testing"
)

func TestUtility_validateBookId(t *testing.T) {
	ctx := context.Background()
	bookId, status := validateBookId(ctx, "TTL")
	if status.IsErr {
		t.Error(status)
	}
	if bookId != "TIT" {
		t.Error(bookId, "should have been revised to TIT")
	}
}
