package timestamp

import (
	"context"
	"dataset"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateYaml() dataset.Status {
	var status dataset.Status
	ctx := context.Background()
	content, err := os.ReadFile("bible_brain_api.json")
	if err != nil {
		return log.Error(ctx, 500, err, "Reading file")
	}
	var bibles []Bible
	err = json.Unmarshal(content, &bibles)
	if err != nil {
		return log.Error(ctx, 500, err, "Error unmarshalling bibles")
	}
	for _, bible := range bibles {
		fmt.Println(bible.Abbr)
		for _, fileset := range bible.Filesets.DBPProd {
			fmt.Println(fileset)
		}
		for _, fileset := range bible.Filesets.DBPVid {
			fmt.Println(fileset)
		}
	}
	return status
}

func GenerateOneYaml(bibleId string, fileset Fileset) dataset.Status {
	var status dataset.Status
	ctx := context.Background()
	filename := filepath.Join("test_data", fileset.FilesetId+".yaml")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return log.Error(ctx, 500, err, "Opening file")
	}
	defer file.Close()

	return status
}

func write(file *os.File, indent int, name string, value string) {
	for i := 0; i < indent; i++ {
		_, _ = file.WriteString(" ")
	}

}

/*
is_new: yes
dataset_name: {FilesetId}
bible_id: {BibleId}
username: GaryNGriswold
email: gary@shortsands.com
output: json
text_data:
bible_brain:
text_usx_edit: yes
audio_data:
bible_brain:
mp3_64: yes
timestamps:
mms_align: yes
testament:
nt: yes
*/
