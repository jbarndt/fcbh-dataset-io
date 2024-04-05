package encode

import (
	"context"
	"fmt"
	"testing"
)

func TestReadDirectory(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	files, status := ReadDirectory(ctx, bibleId, filesetId)
	fmt.Println(status)
	for _, file := range files {
		fmt.Println(file)
	}
}
