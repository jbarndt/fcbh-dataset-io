package diff

import (
	"context"
	"dataset/db"
	"dataset/request"
	"fmt"
	"testing"
)

func TestCompare(t *testing.T) {
	type testCase struct {
		baseDB  string
		project string
		expect  int
	}
	var tests []testCase
	//tests = append(tests, testCase{baseDB: "ATIWBT_SCRIPT.db", project: "ATIWBT_USXEDIT", expect: 2})
	//tests = append(tests, testCase{baseDB: `ENGWEB_WHISPER.db`, project: `ENGWEB_WHISPER_STT`, expect: 90})
	//tests = append(tests, testCase{baseDB: `APFCMU_WHISPER.db`, project: `APFCMU_WHISPER_STT`, expect: 91})
	tests = append(tests, testCase{baseDB: `DYIIBS_WHISPER.db`, project: `DYIIBS_WHISPER_STT`, expect: 91})

	for _, tst := range tests {
		ctx := context.Background()
		user := request.GetTestUser()
		user = ``
		var testament = request.Testament{NT: true}
		var cfg request.CompareSettings
		cfg.LowerCase = true
		cfg.RemovePromptChars = true
		cfg.RemovePunctuation = true
		cfg.DoubleQuotes.Remove = true
		cfg.Apostrophe.Remove = true
		cfg.Hyphen.Remove = true
		cfg.DiacriticalMarks.NormalizeNFD = true
		conn := db.NewDBAdapter(ctx, tst.baseDB)
		compare := NewCompare(ctx, user, tst.project, conn, "eng", testament, cfg)
		filename, status := compare.Process()
		fmt.Println(status, filename)
		if compare.diffCount != tst.expect {
			t.Error(`Expected count is`, tst.expect, ` actual was`, compare.diffCount)
		}
	}
}
