package testing

import (
	"fmt"
	"strings"
	"testing"
)

const PostAudio2WhisperJson = `is_new: no
dataset_name: PostAudioWhisperJson_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 12__post_audio2_whisper.json
audio_data:
  post: {namev4}
text_data:
  speech_to_text:
    whisper:
      model:
        tiny: yes
`

func TestPostAudio2WhisperJsonAPI(t *testing.T) {
	type try struct {
		bibleId  string
		filePath string
		namev4   string
		expected int
	}
	var a try
	a.bibleId = `ENGWEB`
	a.filePath = `ENGWEB/ENGWEBN2DA/B23___02_1John_______ENGWEBN2DA.mp3`
	a.namev4 = `ENGWEBN2DA_B23_1JN_002.mp3`
	destFile := CopyAudio(a.namev4, a.filePath, t)
	a.expected = 72
	var request = strings.Replace(PostAudio2WhisperJson, `{bibleId}`, a.bibleId, 2)
	request = strings.Replace(request, `{namev4}`, destFile, 1)
	stdout, stderr := CurlExec(request, destFile, t)
	fmt.Println(`STDOUT`, stdout)
	fmt.Println(`STDERR`, stderr)
	count := countRecords(stdout)
	if count != a.expected {
		t.Error(`expected,`, a.expected, `found`, count)
	}
}
