package xxxtobedeleted

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TSData struct {
	MediaType    string `json:"media_type"`
	MediaId      string `json:"media_id"`
	PlainText    string `json:"plain_text"`
	ScriptPath   string `json:"script_path"`
	ScriptTSPath string `json:"script_ts_path"`
	LineTSPath   string `json:"line_ts_path"`
	VerseTSPath  string `json:"verse_ts_path"`
	Count        int    `json:"count"`
}

func TestSIL639(t *testing.T) {
	filePath := filepath.Join(os.Getenv(`GOPATH`), `dataset/cli_misc/find_timestamps/TestFilesetList.json`)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	var response []TSData
	err = json.Unmarshal(content, &response)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	for _, ts := range response {
		if ts.MediaType == "audio" || ts.MediaType == "audio_drama" {
			fmt.Println(ts.MediaId)
			fmt.Println(strings.ToLower(ts.MediaId[:3]))
			iso3 := strings.ToLower(ts.MediaId[:3])
			lang, status := FindWhisperCompatibility(ctx, iso3)
			if status.IsErr {
				t.Fatal(status)
			}
			if len(lang) > 0 {
				fmt.Println(ts.MediaId, iso3, lang[0])
			}
		}
	}
}
