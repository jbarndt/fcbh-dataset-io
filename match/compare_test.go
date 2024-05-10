package match

import (
	"context"
	"dataset/db"
	"dataset/fetch"
	"dataset/request"
	"fmt"
	"testing"
)

func TestCompare(t *testing.T) {
	ctx := context.Background()
	user, _ := fetch.GetDBPUser()
	user.Username = ``
	var testament = request.Testament{NT: true}
	var cfg request.CompareSettings
	cfg.LowerCase = true
	cfg.RemovePromptChars = true
	cfg.RemovePunctuation = true
	cfg.DoubleQuotes.Remove = true
	cfg.Apostrophe.Remove = true
	cfg.Hyphen.Remove = true
	cfg.DiacriticalMarks.NormalizeNFD = true
	conn := db.NewDBAdapter(ctx, `ATIWBT_SCRIPT.db`)
	compare := NewCompare(ctx, user, `ATIWBT_USXEDIT`, conn, testament, cfg)
	filename, status := compare.Process()
	fmt.Println(status, filename)
	if compare.diffCount != 2 {
		t.Error(`Expected count is 2, actual was`, compare.diffCount)
	}
}
