package controller

import (
	"dataset_io"
	"dataset_io/fetch"
	"dataset_io/utility"
	"fmt"
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
	var database = utility.NewDBAdapter(databaseName)
	var info = c.fetchMetaData()
	fmt.Println("INFO", info)
	database.InsertIdent(info.BibleId, info.LanguageISO, info.VersionCode, textSource,
		info.LanguageId, info.RolvId, info.Alphabet.Alphabet, info.LanguageName, info.VersionName)
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

}

func (c *Controller) readAudio() {
	// read from acache
}

func (c *Controller) readText() {

}

func (c *Controller) encodeAudio() {

}

func (c *Controller) encodeText() {

}

func (c *Controller) output() {

}
