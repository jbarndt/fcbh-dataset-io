package controller

import (
	"dataset_io"
	"dataset_io/fetch"
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
	var info = c.fetchMetaData()
	fmt.Println("INFO", info)
	// store in audio ident table
	// create DBAdapter

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
