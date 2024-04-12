package controller

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/encode"
	"dataset/fetch"
	log "dataset/logger"
	"dataset/read"
	"dataset/request"
	"dataset/speech_to_text"
	"fmt"
	"strings"
	"time"
)

type Controller struct {
	ctx         context.Context
	yamlRequest []byte
	req         request.Request
	req2        request.Request
	database    db.DBAdapter
}

func NewController(yamlContent []byte) Controller {
	var c Controller
	c.ctx = context.Background()
	c.yamlRequest = yamlContent
	return c
}

func (c *Controller) Process() {
	var status = c.processSteps()
	if status.IsErr {
		c.output(status)
	}
}

func (c *Controller) processSteps() dataset.Status {
	var start = time.Now()
	var status dataset.Status
	// Decode YAML Request File
	reqDecoder := request.NewRequestDecoder(c.ctx)
	c.req, status = reqDecoder.Process(c.yamlRequest)
	if status.IsErr {
		return status
	}
	c.ctx = context.WithValue(context.Background(), `request`, string(c.yamlRequest))
	// Open Database
	dbName := c.req.Required.BibleId + "_" + c.req.TextData.BibleBrain.String() + ".db"
	db.DestroyDatabase(dbName)
	c.database = db.NewDBAdapter(c.ctx, dbName)
	// Fetch Ident Data from DBP
	var info fetch.BibleInfoType
	info, status = c.fetch()
	if status.IsErr {
		return status
	}
	// Read Text Data
	status = c.readText()
	if status.IsErr {
		return status
	}
	if len(info.AudioFilesets) > 0 {
		status = c.speechToText(info.AudioFilesets)
		if status.IsErr {
			return status
		}
	}
	fmt.Println("Duration", time.Since(start))
	return status
}

func (c *Controller) fetch() (fetch.BibleInfoType, dataset.Status) {
	var info fetch.BibleInfoType
	var status dataset.Status
	client := fetch.NewAPIDBPClient(c.ctx, c.req.Required.BibleId)
	info, status = client.BibleInfo()
	if status.IsErr {
		return info, status
	}
	ok := client.FindFilesets(&info, c.req.AudioData.BibleBrain, c.req.TextData.BibleBrain,
		c.req.Testament)
	if ok {
		download := fetch.NewAPIDownloadClient(c.ctx, c.req.Required.BibleId)
		status = download.Download(info)
		if status.IsErr {
			return info, status
		}
	} else {
		var msg = make([]string, 0, 10)
		msg = append(msg, "Requested Fileset is not available")
		for _, rec := range info.DbpProd.Filesets {
			msg = append(msg, fmt.Sprintf("%+v", rec))
		}
		status.IsErr = true
		status.Status = 400
		status.Message = strings.Join(msg, "\n")
		return info, status
	}
	fmt.Println("INFO", info)
	identRec := fetch.CreateIdent(info)
	identRec.TextSource = c.req.TextData.BibleBrain.String()
	c.database.InsertIdent(identRec)
	return info, status
}

func (c *Controller) readText() dataset.Status {
	var status dataset.Status
	bibleId := c.req.Required.BibleId
	if c.req.TextData.BibleBrain.TextUSXEdit {
		status = read.ReadUSXEdit(c.database, bibleId, c.req.Testament)
	} else if c.req.TextData.BibleBrain.TextPlainEdit {
		reader := read.NewDBPTextEditReader(c.req.Required.BibleId, c.database)
		status = reader.Process(c.req.Testament)
	} else if c.req.TextData.BibleBrain.TextPlain {
		reader := read.NewDBPTextReader(c.database)
		status = reader.ProcessDirectory(bibleId, c.req.Testament)
	} else if c.req.TextData.File != `` {
		reader := read.NewScriptReader(c.database)
		var file string
		file, status = reader.FindFile(bibleId)
		if status.IsErr {
			return status
		}
		status = reader.Read(file)
	} else {
		log.Warn(c.ctx, "Could not process ", c.req.TextData)
	}
	if status.IsErr {
		return status
	}
	if c.req.Detail.Words {
		words := read.NewWordParser(c.database)
		status = words.Parse()
	}
	return status
}

func (c *Controller) speechToText(filesets []fetch.FilesetType) dataset.Status {
	var status dataset.Status
	bibleId := c.req.Required.BibleId
	var whisperModel = c.req.TextData.SpeechToText.Whisper.Model.String()
	if whisperModel != `` {
		var whisp = speech_to_text.NewWhisper(bibleId, c.database, whisperModel)
		status = whisp.ProcessDirectory(filesets[0].Id, c.req.Testament)
		if status.IsErr {
			return status
		}
	}
	return status
}

func (c *Controller) timestamps(filesets []fetch.FilesetType) dataset.Status {
	var status dataset.Status
	if c.req.Timestamps.BibleBrain {
		status.IsErr = true
		status.Status = 400
		status.Message = "Timestamps are not yet available"
	} else if c.req.Timestamps.Aeneas {
		aeneas := encode.NewAeneas(c.ctx, c.database, c.req.Required.BibleId, filesets[0].Id)
		status = aeneas.Process(c.req.Required.LanguageISO, c.req.Detail)
	}
	return status
}

func (c *Controller) encodeAudio() dataset.Status {
	var status dataset.Status

	return status
}

func (c *Controller) encodeText() dataset.Status {
	var status dataset.Status

	return status
}

func (c *Controller) matchText() dataset.Status {
	var status dataset.Status

	return status
}

func (c *Controller) output(status dataset.Status) {
	if status.IsErr {
		fmt.Println("IsError:", status.IsErr)
		fmt.Println("Status:", status.Status)
		fmt.Println("GoError:", status.Err)
		fmt.Println("Message:", status.Message)
	}
	//fmt.Println("Response", status)
}
