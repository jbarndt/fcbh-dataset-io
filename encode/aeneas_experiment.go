package encode

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/input"
	"dataset/read"
	"fmt"
	"os"
	"path/filepath"
)

// 1. Load the script into a memory database
// 1b. Load the plain_text_edit into a memory database
// 2. Parsing the script, find those script lines that begin a verse.
// 3. Break the script up into parts that contain 1 or more whole verse
// 4. keep the cue begin and end timestamp with each segment
// 5. for each segment that contains more than 1 verse,

type AeneasExperiment struct {
	ctx context.Context
}

func NewAeneasExperiment(ctx context.Context) AeneasExperiment {
	var a AeneasExperiment
	a.ctx = ctx
	return a
}

func (a *AeneasExperiment) Process() {
	script := a.loadScript(`APFCMUN2DA`)
	fmt.Println(len(script))
}

func (a *AeneasExperiment) loadScript(mediaId string) []db.Script {
	var results []db.Script
	var status dataset.Status
	ts := input.NewTSBucket(a.ctx)
	key := ts.GetKey(input.Script, mediaId, ``, 0)
	fmt.Println(`key:`, key)
	filePath := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), mediaId[:6], mediaId[:8]+`ST.xlsx`)
	ts.DownloadObject(input.TSBucketName, key, filePath)
	fmt.Println(`path`, filePath)
	var conn = db.NewDBAdapter(a.ctx, ":memory:")
	reader := read.NewScriptReader(conn)
	reader.Read(filePath)
	results, status = conn.SelectScripts()
	if status.IsErr {
		panic(status)
	}
	conn.Close()
	fmt.Println(`count:`, len(results))
	return results
}

func (a *AeneasExperiment) loadPlainTextEdit(filePath string) db.DBAdapter {
	var result db.DBAdapter

	return result
}
