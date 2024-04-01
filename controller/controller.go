package controller

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/fetch"
	log "dataset/logger"
	"dataset/read"
	"fmt"
	"os"
	"time"
)

type Controller struct {
	request  dataset.RequestType
	dbName   string
	ctx      context.Context
	database db.DBAdapter
}

func NewController(request dataset.RequestType) *Controller {
	var c Controller
	c.request = request
	c.dbName = request.BibleId + "_" + string(request.TextSource) + ".db"
	db.DestroyDatabase(c.dbName)
	c.ctx = context.WithValue(context.Background(), `request`, request)
	c.database = db.NewDBAdapter(c.ctx, c.dbName)
	return &c
}

func (c *Controller) Process() {
	var start = time.Now()
	var info, ok = c.fetchMetaDataAndFiles()
	if !ok {
		fmt.Println(`Requested Fileset is not available`)
		for _, rec := range info.DbpProd.Filesets {
			fmt.Println(rec)
		}
		os.Exit(0) // Not really, return where
	}
	fmt.Println("INFO", info)
	identRec := fetch.CreateIdent(info)
	identRec.TextSource = string(c.request.TextSource)
	c.database.InsertIdent(identRec)
	c.readText(c.database)
	fmt.Println("Duration", time.Since(start))
}

func (c *Controller) fetchMetaDataAndFiles() (fetch.BibleInfoType, bool) {
	req := c.request
	client := fetch.NewDBPAPIClient(c.ctx, req.BibleId)
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
		log.Warn(c.ctx, "Could not process ", c.request.TextSource)
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
