package request

import (
	"context"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	var d = NewRequestDecoder(context.Background())
	content, err := os.ReadFile(`../controller/cli/request_test.yaml`)
	if err != nil {
		panic(err)
	}
	var req, _ = d.Decode(content)
	req.IsNew = true
	req.BibleId = `EBGESV`
	req.AudioData.File = `file:///where`
	req.AudioData.BibleBrain.MP3_64 = false
	req.AudioData.BibleBrain.OPUS = false
	req.AudioData.POST = `` //`filename`
	req.TextData.NoText = false
	req.TextData.BibleBrain.TextPlain = false
	req.TextData.SpeechToText.Whisper.Model.Medium = true
	req.AudioEncoding.MFCC = true
	req.AudioEncoding.NoEncoding = false
	req.Compare.CompareSettings.Apostrophe.Normalize = true
	req.Compare.CompareSettings.Apostrophe.Remove = false
	status := d.Validate(&req)
	if status.IsErr {
		t.Fatal(status)
	}
}

// I should have a test with multiple error
// I shoud have a test with one selected, not error
// I should have a test with none selected
