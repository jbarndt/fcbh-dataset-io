package controller

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/encode"
	"dataset/fetch"
	"dataset/input"
	log "dataset/logger"
	"dataset/match"
	"dataset/output"
	"dataset/read"
	"dataset/request"
	"dataset/speech_to_text"
	"time"
)

type Controller struct {
	ctx         context.Context
	yamlRequest []byte
	req         request.Request
	user        fetch.DBPUser
	info        fetch.BibleInfoType
	database    db.DBAdapter
}

func NewController(yamlContent []byte) Controller {
	var c Controller
	c.ctx = context.Background()
	c.yamlRequest = yamlContent
	return c
}

func (c *Controller) Process() (string, dataset.Status) {
	var start = time.Now()
	log.SetLevel(log.LOGDEBUG)
	log.SetOutput(c.ctx, `stderr`)
	log.Debug(c.ctx)
	var filename, status = c.processSteps()
	if status.IsErr {
		filename = c.outputStatus(status)
	}
	log.Info(c.ctx, "Duration", time.Since(start))
	log.Debug(c.ctx)
	return filename, status
}

func (c *Controller) processSteps() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	// Decode YAML Request File
	reqDecoder := request.NewRequestDecoder(c.ctx)
	c.req, status = reqDecoder.Process(c.yamlRequest)
	if status.IsErr {
		return filename, status
	}
	var yaml string
	// Stuff YAML request into context
	yaml, status = reqDecoder.Encode(c.req)
	if status.IsErr {
		return filename, status
	}
	c.ctx = context.WithValue(context.Background(), `request`, yaml)
	// Get User
	c.user, status = fetch.GetDBPUser()
	if status.IsErr {
		return filename, status
	}
	// Open Database
	c.database = db.NewerDBAdapter(c.ctx, c.req.Required.IsNew, c.user.Username, c.req.Required.RequestName)
	// Fetch Ident Data from DBP
	c.info, status = c.fetchData()
	if status.IsErr {
		return filename, status
	}
	// Collect Text Input
	var textFiles []input.InputFile
	textFiles, status = c.collectTextInput()
	if status.IsErr {
		return filename, status
	}
	// Collect Audio Input
	var audioFiles []input.InputFile
	audioFiles, status = c.collectAudioInput()
	if status.IsErr {
		return filename, status
	}
	// Read Text Data
	status = c.readText(textFiles)
	if status.IsErr {
		return filename, status
	}
	// Speech to Text
	status = c.speechToText(audioFiles)
	if status.IsErr {
		return filename, status
	}
	// Timestamps
	status = c.timestamps(audioFiles)
	if status.IsErr {
		return filename, status
	}
	// Encode Audio
	status = c.encodeAudio(audioFiles)
	if status.IsErr {
		return filename, status
	}
	// Encode Text
	status = c.encodeText()
	if status.IsErr {
		return filename, status
	}
	// Compare
	if c.req.Compare.BaseProject != `` {
		status = c.matchText()
		if status.IsErr {
			return filename, status
		}
	}
	// Prepare output
	filename, status = c.output()
	return filename, status
}

func (c *Controller) fetchData() (fetch.BibleInfoType, dataset.Status) {
	var info fetch.BibleInfoType
	var status dataset.Status
	client := fetch.NewAPIDBPClient(c.ctx, c.req.Required.BibleId)
	info, status = client.BibleInfo()
	if status.IsErr {
		return info, status
	}
	client.FindFilesets(&info, c.req.AudioData.BibleBrain, c.req.TextData.BibleBrain, c.req.Testament)
	download := fetch.NewAPIDownloadClient(c.ctx, c.req.Required.BibleId)
	status = download.Download(info)
	if status.IsErr {
		return info, status
	}
	//} else {
	//	var msg = make([]string, 0, 10)
	//	msg = append(msg, "Requested Fileset is not available")
	//	for _, rec := range info.DbpProd.Filesets {
	//		msg = append(msg, fmt.Sprintf("%+v", rec))
	//	}
	//	status.IsErr = true
	//	status.Status = 400
	//	status.Message = strings.Join(msg, "\n")
	//	return info, status
	//}
	identRec := client.CreateIdent(info)
	identRec.TextSource = c.req.TextData.BibleBrain.String() // unclear value
	c.database.InsertIdent(identRec)
	return info, status
}

func (c *Controller) collectTextInput() ([]input.InputFile, dataset.Status) {
	var files []input.InputFile
	var status dataset.Status
	var textType string
	if c.req.TextData.BibleBrain.TextUSXEdit {
		textType = `text_usx`
	} else if c.req.TextData.BibleBrain.TextPlainEdit {
		textType = `text_plain`
	} else if c.req.TextData.BibleBrain.TextPlain {
		textType = `text_plain`
	}
	if textType != `` {
		bibleId := c.req.Required.BibleId
		files, status = input.DBPDirectory(c.ctx, bibleId, textType, c.info.TextOTFileset.Id,
			c.info.TextNTFileset.Id, c.req.Testament)
	} else if c.req.TextData.File != `` {
		files, status = input.FileInput(c.ctx, c.req.TextData.File, c.req.Testament)
	} else if c.req.TextData.Http != `` {
		status = log.ErrorNoErr(c.ctx, 400, `Http is not implemented yet`)
	} else if c.req.TextData.AWSS3 != `` {
		files, status = input.AWSS3Input(c.ctx, c.req.TextData.AWSS3, c.req.Testament)
	} else if c.req.TextData.POST {
		status = log.ErrorNoErr(c.ctx, 400, `POST is not implemented yet`)
	}
	return files, status
}

func (c *Controller) collectAudioInput() ([]input.InputFile, dataset.Status) {
	var files []input.InputFile
	var status dataset.Status
	bb := c.req.AudioData.BibleBrain
	if bb.MP3_64 || bb.MP3_16 || bb.OPUS {
		bibleId := c.req.Required.BibleId
		files, status = input.DBPDirectory(c.ctx, bibleId, `audio`, c.info.AudioOTFileset.Id,
			c.info.AudioNTFileset.Id, c.req.Testament)
	} else if c.req.AudioData.File != `` {
		files, status = input.FileInput(c.ctx, c.req.AudioData.File, c.req.Testament)
	} else if c.req.AudioData.Http != `` {
		status = log.ErrorNoErr(c.ctx, 400, `Http is not implemented yet`)
	} else if c.req.AudioData.AWSS3 != `` {
		files, status = input.AWSS3Input(c.ctx, c.req.AudioData.AWSS3, c.req.Testament)
	} else if c.req.AudioData.POST {
		status = log.ErrorNoErr(c.ctx, 400, `POST is not implemented yet`)
	}
	return files, status
}

func (c *Controller) readText(textFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	if len(textFiles) == 0 {
		return status
	}
	//if c.req.TextData.BibleBrain.TextUSXEdit {
	if textFiles[0].MediaType == `text_usx` {
		reader := read.NewUSXParser(c.database)
		status = reader.ProcessFiles(textFiles)
		if status.IsErr {
			return status
		}
	} else if textFiles[0].MediaType == `text_plain` {
		if c.req.TextData.BibleBrain.TextPlainEdit {
			reader := read.NewDBPTextEditReader(c.database, c.req)
			status = reader.Process()
			if status.IsErr {
				return status
			}
		} else { //if c.req.TextData.BibleBrain.TextPlain {
			reader := read.NewDBPTextReader(c.database, c.req.Testament)
			status = reader.ProcessFiles(textFiles)
			if status.IsErr {
				return status
			}
		}
	} else {
		return status // This is not an error, it is nothing to do
	}
	if c.req.Detail.Words {
		words := read.NewWordParser(c.database)
		status = words.Parse()
	}
	return status
}

func (c *Controller) speechToText(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	bibleId := c.req.Required.BibleId
	var whisperModel = c.req.TextData.SpeechToText.Whisper.Model.String()
	if whisperModel != `` {
		var whisper = speech_to_text.NewWhisper(bibleId, c.database, whisperModel)
		status = whisper.ProcessFiles(audioFiles)
		if status.IsErr {
			return status
		}
	}
	return status
}

func (c *Controller) timestamps(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	if c.req.Timestamps.BibleBrain {
		var filesetIds []string
		if c.info.AudioOTFileset.Id != `` {
			filesetIds = append(filesetIds, c.info.AudioOTFileset.Id)
		}
		if c.info.AudioNTFileset.Id != `` {
			filesetIds = append(filesetIds, c.info.AudioNTFileset.Id)
		}
		for _, filesetId := range filesetIds {
			api := fetch.NewAPIDBPTimestamps(c.database, filesetId)
			// Load returns bool, which could be used to invoke aeneas
			_, status = api.LoadTimestamps(c.req.Testament)
			if status.IsErr {
				return status
			}
		}
	} else if c.req.Timestamps.Aeneas {
		bibleId := c.req.Required.BibleId
		aeneas := encode.NewAeneas(c.ctx, c.database, bibleId, c.info.LanguageISO, c.req.Detail)
		status = aeneas.ProcessFiles(audioFiles)
		if status.IsErr {
			return status
		}
	}
	return status
}

func (c *Controller) encodeAudio(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	bibleId := c.req.Required.BibleId
	if c.req.AudioEncoding.MFCC {
		mfcc := encode.NewMFCC(c.ctx, c.database, bibleId, c.req.Detail, 7)
		status = mfcc.ProcessFiles(audioFiles)
		if status.IsErr {
			return status
		}
	}
	return status
}

func (c *Controller) encodeText() dataset.Status {
	var status dataset.Status
	if c.req.TextEncoding.FastText {
		fast := encode.NewFastText(c.ctx, c.database)
		status = fast.Process()
	}
	return status
}

func (c *Controller) matchText() dataset.Status {
	var status dataset.Status
	compare := match.NewCompare(c.ctx, c.req.Compare.BaseProject, c.req.Required.RequestName)
	status = compare.Process()
	return status
}

func (c *Controller) output() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	var out = output.NewOutput(c.ctx, c.database, c.req.Required.RequestName, false, false)
	var records []any
	var meta []output.Meta
	if c.req.Detail.Lines {
		records, meta = out.PrepareScripts()
	} else {
		records, meta = out.PrepareWords()
	}
	if c.req.OutputFormat.CSV {
		filename, status = out.WriteCSV(records, meta)
		if status.IsErr {
			return filename, status
			//return c.outputStatus(status)
		}
	} else if c.req.OutputFormat.JSON {
		filename, status = out.WriteJSON(records, meta)
		if status.IsErr {
			//return c.outputStatus(status)
			return filename, status
		}
	}
	return filename, status
}

func (c *Controller) outputStatus(status dataset.Status) string {
	var filename string
	var status2 dataset.Status
	var out = output.NewOutput(c.ctx, db.DBAdapter{}, c.req.Required.RequestName, false, false)
	if c.req.OutputFormat.CSV {
		filename, status2 = out.CSVStatus(status, true)
	} else if c.req.OutputFormat.JSON {
		filename, status2 = out.JSONStatus(status, true)
	} else {
		filename = status.String()
	}
	if status2.IsErr {
		filename = status2.String()
	}
	return filename
}
