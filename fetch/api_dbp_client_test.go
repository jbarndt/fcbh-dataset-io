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
	req.Testament.OT = true
	ctx := context.Background()
	client := NewAPIDBPClient(ctx, req.Required.BibleId)
	var info, status = client.BibleInfo()
	if status.IsErr {
		t.Error(`Failure in BibleInfo`, status)
	}
	client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
	identRec := client.CreateIdent(info)
	expect := db.Ident{BibleId: `ENGWEB`, AudioOTId: ``, AudioNTId: `ENGWEBN2DA-mp3-64`,
		TextOTId: `ENGWEBO_ET`, TextNTId: `ENGWEBN_ET`, LanguageISO: `eng`, VersionCode: `WEB`, LanguageId: 6414,
		RolvId: 0, Alphabet: `Latn`, LanguageName: `English`, VersionName: `World English Bible`}
	if expect != identRec {
		t.Error("Expected:", expect, "Test:", identRec)
	}
}
