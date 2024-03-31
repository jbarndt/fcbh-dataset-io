package main

import (
	"dataset"
	"dataset/controller"
)

func main() {
	var req dataset.RequestType
	req.BibleId = "BGGWFW"
	//req.BibleId = "ATIWBT"
	//req.TextSource = dataset.USXEDIT
	//req.TextSource = dataset.DBPTEXT
	//req.TextDetail = dataset.LINES
	req.TextDetail = dataset.BOTH
	//req.TextSource = dataset.SCRIPT
	//req.TextSource = dataset.TEXTEDIT
	req.AudioSource = dataset.MP3
	req.Testament = dataset.NT
	var control = controller.NewController(req)
	control.Process()
}
