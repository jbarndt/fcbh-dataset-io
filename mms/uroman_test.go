package mms

import (
	"context"
	"fmt"
	"testing"
)

func TestURoman(t *testing.T) {
	ctx := context.Background()
	var input = []string{"Игорь Стравинский",
		"Игорь",
		"Ντέιβις Καπ",
		"\u0041",
		"123 \u09E6\u09EF \u0966-\u096F \u0660-\u0669", // numbers changed to ascii
		"comma: \u3001 period: \u3002 corners: \u300C\u300D reference: \u203B middle dot: \u30FB",
		"closing double quote: \u030B inverted caret: \u030C upper right: \u031A", // diacriticals are ignored
		"NFC: \u00FC NFD: u\u0308", // NFC processed, NFD diacritical ignored
	}
	results, status := URoman(ctx, "rus", input)
	if status.IsErr {
		t.Fatal(status.String())
	}
	for _, line := range results {
		fmt.Println(line)
	}
}
