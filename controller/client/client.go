package main

import (
	"dataset_io"
	"dataset_io/controller"
)

func main() {
	var req dataset_io.RequestType
	req.BibleId = "ATIWBT"
	req.TextSource = dataset_io.USXEDIT
	var control = controller.NewController(req)
	control.Process()
}
