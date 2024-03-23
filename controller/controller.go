package controller

import (
	"dataset_io"
	"dataset_io/db"
	"dataset_io/fetch"
	"dataset_io/read"
	"fmt"
	"log"
)

type Controller struct {
	request dataset_io.RequestType
}

func NewController(request dataset_io.RequestType) *Controller {
	var c Controller
	c.request = request
	return &c
}

func (c *Controller) Process() {
	textSource := string(c.request.TextSource)
	var databaseName = c.request.BibleId + "_" + textSource + ".db"
	db.DestroyDatabase(databaseName)
	var database = db.NewDBAdapter(databaseName)
	var info = c.fetchMetaData()
	fmt.Println("INFO", info)
	database.InsertIdent(info.BibleId, info.LanguageISO, info.VersionCode, textSource,
		info.LanguageId, info.RolvId, info.Alphabet.Alphabet, info.LanguageName, info.VersionName)
	c.readText(database)
}

func (c *Controller) fetchMetaData() fetch.BibleInfoType {
	req := c.request
	client := fetch.NewDbpApiClient(req.BibleId)
	var info = client.BibleInfo()
	return info
}

func (c *Controller) fetchAudio() {

}

func (c *Controller) fetchText() {
	// Wait on this until talking to Brad. OR should I put in API
}

func (c *Controller) readAudio() {

}

func (c *Controller) readText(database db.DBAdapter) {
	switch c.request.TextSource {
	case dataset_io.USXEDIT:
		read.ReadUSXEdit(database, c.request.BibleId, c.request.Testament)
	default:
		log.Println("Error: Could not process ", c.request.TextSource)
	}
}

func (c *Controller) encodeAudio() {

}

func (c *Controller) encodeText() {

}

func (c *Controller) output() {

}
