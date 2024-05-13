package testing

import (
	"dataset/controller"
	"fmt"
	"strings"
	"testing"
)

const PlainTextEditBBTimestampsScript = `is_new: no
dataset_name: PlainTextEditScript_{bibleId}
bible_id: {bibleId}
audio_data:
  bible_brain:
    mp3_64: yes
timestamps: 
  bible_brain: yes
output_format:
  csv: yes
`

func TestPlainTextBBTimestampsScript(t *testing.T) {
	var bibles = make(map[string]int)
	bibles[`ENGWEB`] = 8250
	//bibles[`ATIWBT`] = 8243
	for bibleId, expected := range bibles {
		var request = strings.Replace(PlainTextEditBBTimestampsScript, `{bibleId}`, bibleId, 2)
		var control = controller.NewController([]byte(request))
		filename, status := control.Process()
		if status.IsErr {
			t.Error(status)
		}
		fmt.Println(filename)
		numLines := NumCVSFileLines(filename, t)
		if numLines != expected {
			t.Error(`Expected `, expected, `records, got`, numLines)
		}
	}
}
