package req

import (
	"testing"
)

func TestValidate(t *testing.T) {
	var req = DecodeFile(`request.yaml`)
	req.Required.BibleId = `EBGESV`
	req.Required.VersionCode = `WBT`
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
	Validate(req)
}

// I should have a test with multiple error
// I shoud have a test with one selected, not error
// I should have a test with none selected
