package main

import (
	"dataset_io"
	"dataset_io/controller"
)

func main() {
	var req dataset_io.RequestType
	req.BibleId = "ATIWBT"
	//req.TextSource = dataset_io.USXEDIT
	//req.TextSource = dataset_io.DBPTEXT
	req.TextDetail = dataset_io.BOTH
	req.TextSource = dataset_io.SCRIPT
	req.Testament = dataset_io.NT
	var control = controller.NewController(req)
	control.Process()
}
