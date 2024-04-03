package fetch

import (
	"context"
	"dataset"
	"dataset/db"
	"testing"
)

func TestAPIDBPClient1(t *testing.T) {
	var req dataset.RequestType
	req.BibleId = `ENGWEB`
	req.AudioSource = dataset.MP3
	req.TextSource = dataset.TEXTEDIT
	req.Testament = dataset.NT
	ctx := context.Background()
	client := NewAPIDBPClient(ctx, req.BibleId)
	var info, status = client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioSource, req.TextSource, req.Testament)
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
