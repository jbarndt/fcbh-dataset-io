package main

import (
	"dataset_io"
	"dataset_io/controller"
)

func main() {
	var req dataset_io.RequestType
	req.BibleId = "BGGWFW"
	//req.TextSource = dataset_io.USXEDIT
	//req.TextSource = dataset_io.DBPTEXT
	req.TextDetail = dataset_io.BOTH
	req.TextSource = dataset_io.SCRIPT
	//req.TextSource = dataset_io.TEXTEDIT
	req.Testament = dataset_io.NT
	var control = controller.NewController(req)
	control.Process()
}
