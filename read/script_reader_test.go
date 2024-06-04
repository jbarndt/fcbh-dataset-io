package read

import (
	"context"
	"dataset/db"
	"dataset/input"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestScriptReader(t *testing.T) {
	bibleId := `ATIWBT`
	database := bibleId + "_SCRIPT.db"
	db.DestroyDatabase(database)
	conn := db.NewDBAdapter(context.Background(), database)
	script := NewScriptReader(conn)
	filename := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, bibleId+`N2ST.xlsx`)
	fmt.Println(`Filename:`, filename)
	status := script.Read(filename)
	if status.IsErr {
		t.Fatal(status)
	}
	//}
	conn.Close()
}

// TestScriptHeaders is to investigate the format of script columns for position of items
func TestScriptHeaders(t *testing.T) {
	ctx := context.Background()
	ts := input.NewTSBucket(ctx)
	list := ts.ListPrefix(input.TSBucketName, input.LatinN2)
	for count, item := range list {
		if count > 400 {
			break
		}
		//fmt.Println(item)
		objs := ts.ListObjects(input.TSBucketName, item+input.Script)
		key := objs[0]
		//fmt.Println(key)
		filename := filepath.Base(key)
		filePath := filepath.Join(os.Getenv(`FCBH_DATASET_TMP`), filename)
		//fmt.Println(filePath)
		ts.DownloadObject(input.TSBucketName, key, filePath)
		file, err := excelize.OpenFile(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		sheets := file.GetSheetList()
		rows, err := file.GetRows(sheets[0])
		if err != nil {
			t.Fatal(err)
		}
		for col, cell := range rows[0] {
			fmt.Print(col, `:`, cell, ` `)
		}
		fmt.Print("\n")
		file.Close()
		os.Remove(filePath)
	}
}
