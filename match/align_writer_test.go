package match

import (
	"context"
	"dataset/db"
	"dataset/generic"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAlignWriter(t *testing.T) {
	ctx := context.Background()
	//var dataset = "N2YPM_JMD"
	//var dataset = "ENGWEB_align"
	var dataset = "ENGWEB_align_mp3"
	//dbPath := filepath.Join(os.Getenv("HOME"), "FCBH2024", "GaryNTest", "PlainTextEditScript_ENGWEB.db")
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".db")
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn, "eng", "")
	audioDir := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), "ENGWEB", "ENGWEBN2DA-mp3-64")
	faVerses, filenameMap, status := calc.Process(audioDir)
	if status.IsErr {
		t.Fatal(status)
	}
	fmt.Println(len(faVerses), len(filenameMap))
	writer := NewAlignWriter(ctx)
	filename, status := writer.WriteReport(dataset, faVerses, filenameMap)
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}

func TestAlignWriter_JsonInput(t *testing.T) {
	ctx := context.Background()
	var dataset = "ENGWEB_align_mp3"
	dbPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".db")
	conn := db.NewDBAdapter(ctx, dbPath)
	calc := NewAlignErrorCalc(ctx, conn, "eng", "")
	filenameMap, status2 := calc.generateBookChapterFilenameMap()
	if status2.IsErr {
		t.Fatal(status2)
	}
	bytes, err := os.ReadFile("mms_asr_align.json")
	if err != nil {
		t.Fatal(err)
	}
	var lines []generic.AlignLine
	err = json.Unmarshal(bytes, &lines)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	for _, line := range lines {
		for _, ch := range line.Chars {
			if ch.IsASR {
				count++
			}
		}
	}
	fmt.Println(count)
	writer := NewAlignWriter(ctx)
	filename, status3 := writer.WriteReport(dataset, lines, filenameMap)
	if status3.IsErr {
		t.Fatal(status3)
	}
	fmt.Println("Report Filename", filename)
	revisedName := filepath.Join(os.Getenv("GOPROJ"), "dataset", "match", dataset+".html")
	_ = os.Rename(filename, revisedName)
	fmt.Println("Report Filename", revisedName)
}
