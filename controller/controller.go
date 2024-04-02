package controller

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/fetch"
	log "dataset/logger"
	"dataset/read"
	"fmt"
	"strings"
	"time"
)

type Controller struct {
	request  dataset.RequestType
	dbName   string
	ctx      context.Context
	database db.DBAdapter
}

func NewController(request dataset.RequestType) Controller {
	var c Controller
	c.request = request
	c.dbName = request.BibleId + "_" + string(request.TextSource) + ".db"
	db.DestroyDatabase(c.dbName)
	c.ctx = context.WithValue(context.Background(), `request`, request)
	c.database = db.NewDBAdapter(c.ctx, c.dbName)
	return c
}

func (c *Controller) Process() {
	var start = time.Now()
	var info, status, ok = c.fetchMetaDataAndFiles()
	if !status.IsErr {
		if !ok {
			var msg = make([]string, 0, 10)
			msg = append(msg, "Requested Fileset is not available")
			for _, rec := range info.DbpProd.Filesets {
				msg = append(msg, fmt.Sprintf("%+v", rec))
			}
			status.Message = strings.Join(msg, "\n")
			c.output(status)
		}
		fmt.Println("INFO", info)
		identRec := fetch.CreateIdent(info)
		identRec.TextSource = string(c.request.TextSource)
		c.database.InsertIdent(identRec)
		c.readText(c.database)
	}
	fmt.Println("Duration", time.Since(start))
	c.output(status)
}

func (c *Controller) fetchMetaDataAndFiles() (fetch.BibleInfoType, dataset.Status, bool) {
	var info fetch.BibleInfoType
	var status dataset.Status
	var ok bool
	req := c.request
	client := fetch.NewAPIDBPClient(c.ctx, req.BibleId)
	info, status = client.BibleInfo()
	if !status.IsErr {
		ok = client.FindFilesets(&info, req.AudioSource, req.TextSource, req.Testament)
		if ok {
			download := fetch.NewAPIDownloadClient(c.ctx, req.BibleId)
			status = download.Download(info)
		}
	}
	return info, status, ok
}

func (c *Controller) readText(database db.DBAdapter) dataset.Status {
	var status dataset.Status
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
		var file string
		file, status = reader.FindFile(c.request.BibleId)
		if status.IsErr {
			return status
		}
		reader.Read(file)
	case dataset.NOTEXT:
	default:
		log.Warn(c.ctx, "Could not process ", c.request.TextSource)
	}
	if c.request.TextDetail == dataset.WORDS || c.request.TextDetail == dataset.BOTH {
		words := read.NewWordParser(database)
		words.Parse()
	}
	return status
}

func (c *Controller) encodeAudio() {

}

func (c *Controller) encodeText() {

}

func (c *Controller) output(status dataset.Status) {
	if status.IsErr {
		fmt.Println("IsError", status.IsErr)
		fmt.Println("Status", status.Status)
		fmt.Println("Error", status.Err)
		fmt.Println("Message", status.Message)
	}
	fmt.Println("Response", status)
}
