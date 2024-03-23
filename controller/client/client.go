package main

import (
	"dataset_io"
	"dataset_io/controller"
)

func main() {
	var req dataset_io.RequestType
	req.BibleId = "ATIWBT"
	var control = controller.NewController(req)
	control.Process()
}
