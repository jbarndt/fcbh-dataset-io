package controller

import (
	"dataset"
	"dataset/db"
	"dataset/fetch"
	"dataset/read"
	"fmt"
	"log"
	"os"
)

type Controller struct {
	request dataset.RequestType
}

func NewController(request dataset.RequestType) *Controller {
	var c Controller
	c.request = request
	return &c
}

func (c *Controller) Process() {
	textSource := string(c.request.TextSource)
	var databaseName = c.request.BibleId + "_" + textSource + ".db"
	var info, ok = c.fetchMetaDataAndFiles()
	if !ok {
		fmt.Println(`Requested Fileset is not available`)
		for _, rec := range info.DbpProd.Filesets {
			fmt.Println(rec)
		}
		os.Exit(0) // Not really, return where
	}
	fmt.Println("INFO", info)
	db.DestroyDatabase(databaseName)
	var database = db.NewDBAdapter(databaseName)
	audioFSId := fetch.ConcatFilesetId(info.AudioFilesets)
	textFSId := fetch.ConcatFilesetId(info.TextFilesets)
	database.InsertIdent(info.BibleId, audioFSId, textFSId, info.LanguageISO, info.VersionCode, textSource,
		info.LanguageId, info.RolvId, info.Alphabet.Alphabet, info.LanguageName, info.VersionName)
	c.readText(database)
}

func (c *Controller) fetchMetaDataAndFiles() (fetch.BibleInfoType, bool) {
	req := c.request
	client := fetch.NewDBPAPIClient(req.BibleId)
	var info = client.BibleInfo()
	ok := client.FindFilesets(&info, req.AudioSource, req.TextSource, req.Testament)
	if ok {
		client.Download(info)
	}
	return info, ok
}

func (c *Controller) readText(database db.DBAdapter) {
	switch c.request.TextSource {
	case dataset.USXEDIT:
		read.ReadUSXEdit(database, c.request.BibleId, c.request.Testament)
	case dataset.TEXTEDIT:
		reader := read.NewDBPTextEditReader(c.request.BibleId, database)
		reader.Process(c.request.Testament)
	case dataset.DBPTEXT:
		reader := read.NewDBPTextReader(database)
		reader.ProcessDirectory(c.request.BibleId, c.request.Testament)
	case dataset.SCRIPT:
		reader := read.NewScriptReader(database)
		file := reader.FindFile(c.request.BibleId)
		reader.Read(file)
	default:
		log.Println("Error: Could not process ", c.request.TextSource)
	}
	if c.request.TextDetail == dataset.WORDS || c.request.TextDetail == dataset.BOTH {
		words := read.NewWordParser(database)
		words.Parse()
	}
}

func (c *Controller) encodeAudio() {

}

func (c *Controller) encodeText() {

}

func (c *Controller) output() {

}
