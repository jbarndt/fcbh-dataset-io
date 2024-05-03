package request

import (
	"context"
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	var d = NewRequestDecoder(context.Background())
	content, err := os.ReadFile(`../controller/client/request_test.yaml`)
	if err != nil {
		panic(err)
	}
	var req, _ = d.Decode(content)
	req.IsNew = true
	req.BibleId = `EBGESV`
	req.AudioData.File = `file:///where`
	req.AudioData.BibleBrain.MP3_64 = true
	req.AudioData.BibleBrain.OPUS = true
	req.AudioData.POST = true
	req.TextData.NoText = true
	req.TextData.BibleBrain.TextPlain = true
	req.TextData.SpeechToText.Whisper.Model.Medium = true
	req.AudioEncoding.MFCC = true
	req.AudioEncoding.NoEncoding = true
	req.Compare.CompareSettings.Apostrophe.Normalize = true
	req.Compare.CompareSettings.Apostrophe.Remove = true
	d.Validate(&req)
}

// I should have a test with multiple error
// I shoud have a test with one selected, not error
// I should have a test with none selected
