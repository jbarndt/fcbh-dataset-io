package tests

import (
	"context"
	"dataset/db"
	"dataset/generic"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"regexp"
	"strings"
	"testing"
)

const usxVsPlain = `is_new: yes
dataset_name: 3b_usx_vs_plain_{BIBLE_ID}
bible_id: {BIBLE_ID}
username: GaryNTest
email: gary@shortsands.com
output_file: 3b_{TEXT_TYPE}_{BIBLE_ID}.sqlite
text_data:
  bible_brain:
    {TEXT_TYPE}: yes
detail:
  words: yes
`

func TestTextReadDirect(t *testing.T) {
	var tests1 []SqliteTest
	tests1 = append(tests1, SqliteTest{"SELECT count(*) FROM scripts", 8213})
	tests1 = append(tests1, SqliteTest{"SELECT count(*) FROM words", 381444})
	test1 := strings.ReplaceAll(usxVsPlain, "{BIBLE_ID}", "ENGWEB")
	test1USX := strings.ReplaceAll(test1, "{TEXT_TYPE}", "text_usx_edit")
	database1 := DirectSqlTest(test1USX, tests1, t)
	var tests2 []SqliteTest
	tests2 = append(tests2, SqliteTest{"SELECT count(*) FROM scripts", 8218})
	tests2 = append(tests2, SqliteTest{"SELECT count(*) FROM words", 381444})
	text2TXT := strings.ReplaceAll(test1, "{TEXT_TYPE}", "text_plain_edit")
	database2 := DirectSqlTest(text2TXT, tests2, t)
	diffCount := DifferenceTest(database1, database2)
	if diffCount != 53 {
		t.Error("DiffCount is expected to be 53")
	}
}

func DifferenceTest(database1 string, database2 string) int {
	os.Setenv("FORCE_COLOR", "1")
	ctx := context.Background()
	diffMatch := diffmatchpatch.New()
	conn1 := db.NewDBAdapter(ctx, "./"+database1)
	records1, status := conn1.SelectScripts()
	if status.IsErr {
		panic(status)
	}
	var usxMap = make(map[string]string)
	for _, rec := range records1 {
		var lf generic.LineRef
		lf.BookId = rec.BookId
		lf.ChapterNum = rec.ChapterNum
		lf.VerseStr = rec.VerseStr
		usxMap[lf.ComposeKey()] = rec.ScriptText
	}
	conn1.Close()
	conn2 := db.NewDBAdapter(ctx, "./"+database2)
	plainRec2, _ := conn2.SelectScripts()
	var diffCount = 0
	for _, rec := range plainRec2 {
		var lf generic.LineRef
		lf.BookId = rec.BookId
		lf.ChapterNum = rec.ChapterNum
		lf.VerseStr = rec.VerseStr
		lineRef := lf.ComposeKey()
		usxTxt, ok := usxMap[lineRef]
		if !ok {
			usxTxt = ""
		}
		usxTxt = strings.TrimSpace(usxTxt)
		plainTxt := rec.ScriptText
		re := regexp.MustCompile(`\n\s+`)
		plainTxt = re.ReplaceAllString(plainTxt, " ")
		diffs := diffMatch.DiffMain(usxTxt, plainTxt, false)
		diffs = diffMatch.DiffCleanupMerge(diffs)
		if len(diffs) > 1 || len(diffs) > 0 && diffs[0].Type != diffmatchpatch.DiffEqual {
			diffCount++
			fmt.Println(lineRef, "usxTxt", usxTxt)
			fmt.Println(lineRef, "plnTxt", plainTxt)
			fmt.Println(lineRef, diffMatch.DiffPrettyText(diffs))
			fmt.Println(lineRef, diffs)
			fmt.Println()
		}
	}
	fmt.Println("DiffCount", diffCount)
	return diffCount
}
