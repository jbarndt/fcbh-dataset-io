package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

type CtlTest struct {
	BibleId   string
	Name      string
	TextNtId  string
	AudioNTId string
	Language  string
	Expected  int
}

func DirectTestUtility(requestYaml string, tests []CtlTest, t *testing.T) {
	ctx := context.Background()
	for _, tst := range tests {
		var req = strings.Replace(requestYaml, `{bibleId}`, tst.BibleId, 2)
		var control = controller.NewController(ctx, []byte(req))
		filename, status := control.Process()
		if status.IsErr {
			t.Fatal(status)
		}
		fmt.Println(filename)
		numLines := NumFileLines(filename, t)
		if numLines != tst.Expected {
			t.Error(`Expected `, tst.Expected, `records, got`, numLines)
		}
		var decoder = request.NewRequestDecoder(ctx)
		reqObj, status := decoder.Decode([]byte(req))
		if status.IsErr {
			t.Fatal(status)
		}
		identTest(reqObj.DatasetName, t, request.TextPlainEdit, ``,
			tst.TextNtId, ``, tst.AudioNTId, tst.Language)
	}
}
