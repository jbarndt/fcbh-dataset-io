package testing

import (
	"context"
	"dataset/controller"
	"dataset/request"
	"fmt"
	"strings"
	"testing"
)

const AudioWhisperJson = `is_new: yes
dataset_name: AudioWhisperJson_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 09__audio_whisper_json.json
testament: # Choose one or both
  nt_books: ["TIT"]
audio_data:
  bible_brain:
    mp3_64: yes
text_data:
  speech_to_text:
    whisper:
      model:
        tiny: yes
`

func TestAudioWhisperJson(t *testing.T) {
	var bibles = make(map[string]int)
	bibles[`ENGWEB`] = 68
	ctx := context.Background()
	for bibleId, expected := range bibles {
		var req = strings.Replace(AudioWhisperJson, `{bibleId}`, bibleId, 2)
		ctrl := controller.NewController(ctx, []byte(req))
		filename, status := ctrl.Process()
		fmt.Println("Filename", filename, status)
		if status.IsErr {
			t.Fatal(status)
		}
		numLines := NumJSONFileLines(filename, t)
		if numLines != expected {
			t.Error(`Expected `, expected, `records, got`, numLines)
		}
		identTest(`AudioWhisperJson_`+bibleId, t, request.TextSTT, ``,
			`ENGWEBN_TT`, ``, `ENGWEBN2DA`, `eng`)
	}
}

func TestAudioWhisperJsonCLI(t *testing.T) {
	var bibles = make(map[string]int)
	bibles[`ENGWEB`] = 68
	for bibleId, expected := range bibles {
		var request = strings.Replace(AudioWhisperJson, `{bibleId}`, bibleId, 2)
		stdout, stderr := CLIExec(request, t)
		fmt.Println(`STDOUT:`, stdout)
		fmt.Println(`STDERR:`, stderr)
		filename := ExtractFilenaame(request)
		numLines := NumJSONFileLines(filename, t)
		if numLines != expected {
			t.Error(`Expected `, expected, `records, got`, numLines)
		}
	}
}
