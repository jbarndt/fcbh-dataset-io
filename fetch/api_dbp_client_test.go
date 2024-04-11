package fetch

import (
	"context"
	"dataset/db"
	"dataset/request"
	"testing"
)

func TestAPIDBPClient1(t *testing.T) {
	var req request.Request
	req.Required.BibleId = `ENGWEB`
	req.AudioData.BibleBrain.MP3_64 = true       // = dataset.MP3
	req.TextData.BibleBrain.TextPlainEdit = true // = dataset.TEXTEDIT
	req.Testament.NT = true                      // = dataset.NT
	ctx := context.Background()
	client := NewAPIDBPClient(ctx, req.Required.BibleId)
	var info, status = client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
		if ok {
			identRec := CreateIdent(info)
			expect := db.Ident{BibleId: `ENGWEB`, AudioFilesetId: `ENGWEBN2DA`,
				TextFilesetId: `ENGWEBN_ET`, LanguageISO: `eng`, VersionCode: `WEB`, LanguageId: 6414,
				RolvId: 0, Alphabet: `Latn`, LanguageName: `English`, VersionName: `World English Bible`}
			if expect != identRec {
				t.Error("Expected:", expect, "Test:", identRec)
			}
		} else {
			t.Error(`Find Filesets not ok`)
		}
	} else {
		t.Error(`API Error`)
	}
}
