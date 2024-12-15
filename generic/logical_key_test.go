package generic

import (
	"testing"
)

func TestReference_ByKey(t *testing.T) {
	var a = RefByKey("NUM 22:12")
	if a.BookId() != "NUM" {
		t.Error("BookId should be NUM")
	}
	if a.ChapterNum() != 22 {
		t.Error("ChapterNum should be 22")
	}
	if a.VerseStr() != "12" {
		t.Error("VerseStr should be 12")
	}
}

func TestReference_ByParts(t *testing.T) {
	var a = RefByParts("NUM", 22, "12", "", 0)
	if a.LogicalKey() != "NUM 22:12" {
		t.Error("BookId should be NUM 22:12")
	}
}
