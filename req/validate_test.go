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
	req.AudioData.POST = true
	Validate(req)
}
