package fetch

import (
	"dataset"
	"dataset/db"
	"fmt"
	"os"
	"testing"
)

func TestDBPAPIClient(t *testing.T) {
	var bibleId = `ENGWEB`
	var textSource = `DBPTEXT`
	var databaseName = bibleId + `_` + textSource + `.db`
	var info, ok = fetchMetaDataAndFiles(bibleId)
	if !ok {
		fmt.Println(`Requested Fileset is not available`)
		for _, rec := range info.DbpProd.Filesets {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", rec.Id, rec.Type, rec.Size, rec.Codec, rec.Bitrate)
		}
		os.Exit(0) // Not really, return where
	}
	//fmt.Println("INFO", info)
	db.DestroyDatabase(databaseName)
	var database = db.NewDBAdapter(databaseName)
	identRec := CreateIdent(info)
	identRec.TextSource = textSource
	database.InsertIdent(identRec)
}

func fetchMetaDataAndFiles(bibleId string) (BibleInfoType, bool) {
	var req dataset.RequestType
	req.BibleId = bibleId
	req.AudioSource = dataset.MP3
	req.TextSource = dataset.TEXTEDIT
	req.Testament = dataset.NT
	client := NewDBPAPIClient(req.BibleId)
	var info = client.BibleInfo()
	ok := client.FindFilesets(&info, req.AudioSource, req.TextSource, req.Testament)
	if ok {
		client.Download(info)
	}
	return info, ok
}
