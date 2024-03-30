package main

import (
	"dataset"
	"dataset/controller"
)

func main() {
	var req dataset.RequestType
	req.BibleId = "BGGWFW"
	req.TextSource = dataset.USXEDIT
	//req.TextSource = dataset.DBPTEXT
	req.TextDetail = dataset.BOTH
	//req.TextSource = dataset.SCRIPT
	//req.TextSource = dataset_io.TEXTEDIT
	req.Testament = dataset.NT
	var control = controller.NewController(req)
	control.Process()
}
