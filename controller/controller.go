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
	"dataset/mms"
	"dataset/output"
	"dataset/read"
	"dataset/request"
	"dataset/run_control"
	"dataset/speech_to_text"
	"dataset/timestamp"
	"time"
)

type Controller struct {
	ctx         context.Context
	yamlRequest []byte
	req         request.Request
	bucket      run_control.RunBucket
	user        fetch.DBPUser
	ident       db.Ident
	database    db.DBAdapter
	postFiles   *input.PostFiles
}

func NewController(ctx context.Context, yamlContent []byte) Controller {
	var c Controller
	c.ctx = ctx
	c.yamlRequest = yamlContent
	c.bucket = run_control.NewRunBucket(ctx, yamlContent)
	c.bucket.IsUnitTest = false // set to true when testing to make RunBucket work.
	return c
}

func (c *Controller) SetPostFiles(postFiles *input.PostFiles) {
	c.postFiles = postFiles
}

func (c *Controller) Process() (string, dataset.Status) {
	var start = time.Now()
	if c.postFiles != nil {
		defer c.postFiles.RemoveDir()
	}
	log.Debug(c.ctx)
	var filename, status = c.processSteps()
	if status.IsErr {
		filename = c.outputStatus(status)
	}
	c.bucket.AddOutput(filename)
	log.Info(c.ctx, "Duration", time.Since(start))
	log.Debug(c.ctx)
	c.bucket.PersistToBucket()
	return filename, status
}

func (c *Controller) processSteps() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	// Decode YAML Request File
	log.Info(c.ctx, "Parse .yaml file.")
	reqDecoder := request.NewRequestDecoder(c.ctx)
	c.req, status = reqDecoder.Process(c.yamlRequest)
	if status.IsErr {
		return filename, status
	}
	c.ctx = context.WithValue(c.ctx, `request`, string(c.yamlRequest))
	// Get User
	log.Info(c.ctx, "Fetch Bible Brain data.")
	c.user, status = fetch.GetDBPUser(c.req)
	if status.IsErr {
		return filename, status
	}
	// Open Database
	c.database, status = db.NewerDBAdapter(c.ctx, c.req.IsNew, c.user.Username, c.req.DatasetName)
	if status.IsErr {
		return filename, status
	}
	defer c.database.Close()
	c.bucket.AddDatabase(c.database)
	// Fetch Ident Data from Ident
	c.ident, status = c.database.SelectIdent()
	if status.IsErr {
		return filename, status
	}
	// Update Ident Data from DBP
	c.ident, status = c.fetchData()
	if status.IsErr {
		if c.req.TextData.AnyBibleBrain() || c.req.AudioData.AnyBibleBrain() {
			return filename, status
		}
	}
	// Collect Text Input
	var textFiles []input.InputFile
	if !c.req.TextData.NoText {
		log.Info(c.ctx, "Load text files.")
		textFiles, status = c.collectTextInput()
		if status.IsErr {
			return filename, status
		}
	}
	// Collect Audio Input
	var audioFiles []input.InputFile
	if !c.req.AudioData.NoAudio {
		log.Info(c.ctx, "Load audio files.")
		audioFiles, status = c.collectAudioInput()
		if status.IsErr {
			return filename, status
		}
	}
	// Update Ident Table
	status = input.UpdateIdent(c.database, &c.ident, textFiles, audioFiles)
	if status.IsErr {
		return filename, status
	}
	// Read Text Data
	if !c.req.TextData.NoText {
		log.Info(c.ctx, "Read and parse text files.")
		status = c.readText(textFiles)
		if status.IsErr {
			return filename, status
		}
	}
	// Timestamps
	if !c.req.Timestamps.NoTimestamps {
		log.Info(c.ctx, "Read or create audio timestamp data.")
		status = c.timestamps(audioFiles)
		if status.IsErr {
			return filename, status
		}
	}
	// Copy for STT
	if !c.req.TextData.NoText &&
		//c.req.Compare.BaseDataset == `` &&
		!c.req.SpeechToText.NoSpeechToText {
		c.req.Compare.BaseDataset = c.database.Project
		// This makes a copy of database, and closes it.  Names the new database *_audio, and returns new
		c.database, status = c.database.CopyDatabase(`_audio`)
		if status.IsErr {
			return filename, status
		}
		c.bucket.AddDatabase(c.database)
		status = c.database.UpdateEraseScriptText()
		if status.IsErr {
			return filename, status
		}
	}
	// Speech to Text
	if !c.req.SpeechToText.NoSpeechToText {
		log.Info(c.ctx, "Perform speech to text.")
		status = c.speechToText(audioFiles)
		if status.IsErr {
			return filename, status
		}
	}
	// Encode Audio
	if !c.req.AudioEncoding.NoEncoding {
		log.Info(c.ctx, "Perform audio encoding.")
		status = c.encodeAudio(audioFiles)
		if status.IsErr {
			return filename, status
		}
	}
	// Encode Text
	if !c.req.TextEncoding.NoEncoding {
		log.Info(c.ctx, "Perform text encoding.")
		status = c.encodeText()
		if status.IsErr {
			return filename, status
		}
	}
	// Compare
	if c.req.OutputFormat.HTML {
		log.Info(c.ctx, "Perform text comparison.")
		filename, status = c.matchText()
		return filename, status // return whether success or not
	}
	// Prepare output
	log.Info(c.ctx, "Generate output.")
	if c.req.OutputFormat.Sqlite {
		filename = c.database.DatabasePath
	} else {
		filename, status = c.output()
	}
	return filename, status
}

func (c *Controller) fetchData() (db.Ident, dataset.Status) {
	var status dataset.Status
	var info fetch.BibleInfoType
	client := fetch.NewAPIDBPClient(c.ctx, c.req.BibleId)
	info, status = client.BibleInfo()
	if status.IsErr {
		return c.ident, status
	}
	client.FindFilesets(&info, c.req.AudioData.BibleBrain, c.req.TextData.BibleBrain, c.req.Testament)
	download := fetch.NewAPIDownloadClient(c.ctx, c.req.BibleId, c.req.Testament)
	status = download.Download(info)
	if status.IsErr {
		return c.ident, status
	}
	c.ident = client.UpdateIdent(c.ident, info, c.req.TextData.BibleBrain.TextType())
	return c.ident, status
}

func (c *Controller) collectTextInput() ([]input.InputFile, dataset.Status) {
	var files []input.InputFile
	var status dataset.Status
	var expectFiles = true
	bb := c.req.TextData.BibleBrain
	if bb.TextPlain || bb.TextPlainEdit || bb.TextUSXEdit {
		files, status = input.DBPDirectory(c.ctx, c.req.BibleId, c.ident.TextSource, c.ident.TextOTId,
			c.ident.TextNTId, c.req.Testament)
	} else if c.req.TextData.File != `` {
		files, status = input.FileInput(c.ctx, c.req.TextData.File, c.req.Testament)
	} else if c.req.TextData.AWSS3 != `` {
		files, status = input.AWSS3Input(c.ctx, c.req.TextData.AWSS3, c.req.Testament)
	} else if c.req.TextData.POST != `` && c.postFiles != nil {
		files, status = c.postFiles.PostInput("text", c.req.Testament)
	} else {
		expectFiles = false
	}
	if expectFiles && len(files) == 0 {
		status = log.ErrorNoErr(c.ctx, 400, `No text files found for`, c.ident.TextSource)
	}
	return files, status
}

func (c *Controller) collectAudioInput() ([]input.InputFile, dataset.Status) {
	var files []input.InputFile
	var status dataset.Status
	var expectFiles = true
	bb := c.req.AudioData.BibleBrain
	if bb.MP3_64 || bb.MP3_16 || bb.OPUS {
		bibleId := c.req.BibleId
		files, status = input.DBPDirectory(c.ctx, bibleId, request.Audio, c.ident.AudioOTId,
			c.ident.AudioNTId, c.req.Testament)
	} else if c.req.AudioData.File != `` {
		files, status = input.FileInput(c.ctx, c.req.AudioData.File, c.req.Testament)
	} else if c.req.AudioData.AWSS3 != `` {
		files, status = input.AWSS3Input(c.ctx, c.req.AudioData.AWSS3, c.req.Testament)
	} else if c.req.AudioData.POST != `` && c.postFiles != nil {
		files, status = c.postFiles.PostInput("audio", c.req.Testament)
	} else {
		expectFiles = false
	}
	if expectFiles && len(files) == 0 {
		status = log.ErrorNoErr(c.ctx, 400, `No audio files found for`, c.ident.AudioNTId)
	}
	return files, status
}

func (c *Controller) readText(textFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	if len(textFiles) == 0 {
		return status
	}
	if textFiles[0].MediaType == request.TextUSXEdit {
		reader := read.NewUSXParser(c.database)
		status = reader.ProcessFiles(textFiles)
		if status.IsErr {
			return status
		}
	} else if textFiles[0].MediaType == request.TextPlainEdit {
		reader := read.NewDBPTextEditReader(c.database, c.req)
		status = reader.Process()
		if status.IsErr {
			return status
		}
	} else if textFiles[0].MediaType == request.TextPlain {
		reader := read.NewDBPTextReader(c.database, c.req.Testament)
		status = reader.ProcessFiles(textFiles)
		if status.IsErr {
			return status
		}
	} else if textFiles[0].MediaType == request.TextScript {
		reader := read.NewScriptReader(c.database, c.req.Testament)
		status = reader.ProcessFiles(textFiles)
		if status.IsErr {
			return status
		}
	} else if textFiles[0].MediaType == request.TextCSV {
		reader := read.NewCSVReader(c.database)
		status = reader.ProcessFiles(textFiles)
		if status.IsErr {
			return status
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

func (c *Controller) timestamps(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	if c.req.Timestamps.BibleBrain {
		//
		// Why isn't Bible Brain just processing input??
		//
		var filesetIds = []string{c.ident.AudioOTId, c.ident.AudioNTId}
		for _, filesetId := range filesetIds {
			if filesetId != `` {
				api := fetch.NewAPIDBPTimestamps(c.database, filesetId)
				// Load returns bool, which could be used to invoke aeneas
				_, status = api.LoadTimestamps(c.req.Testament)
				if status.IsErr {
					return status
				}
			}
		}
	} else if c.req.Timestamps.Aeneas {
		bibleId := c.req.BibleId
		aeneas := encode.NewAeneas(c.ctx, c.database, bibleId, c.ident.LanguageISO, c.req.Detail)
		status = aeneas.ProcessFiles(audioFiles)
	} else if c.req.Timestamps.TSBucket {
		var ts timestamp.TSBucket
		ts, status = timestamp.NewTSBucket(c.ctx, c.database)
		if !status.IsErr {
			status = ts.ProcessFiles(audioFiles)
		}
	} else if c.req.Timestamps.MMSFAVerse {
		var ts mms.ForcedAlign
		ts = mms.NewForcedAlign(c.ctx, c.database, c.ident.LanguageISO, c.req.AltLanguage)
		status = ts.ProcessFiles(audioFiles)
	} else if c.req.Timestamps.MMSAlign {
		var ts mms.MMSFA
		ts = mms.NewMMSFA(c.ctx, c.database, c.ident.LanguageISO, c.req.AltLanguage)
		status = ts.ProcessFiles(audioFiles)
	}
	return status
}

func (c *Controller) speechToText(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	bibleId := c.req.BibleId
	if c.req.SpeechToText.MMS {
		var asr mms.MMSASR
		asr = mms.NewMMSASR(c.ctx, c.database, c.ident.LanguageISO, c.req.AltLanguage)
		status = asr.ProcessFiles(audioFiles)
	} else {
		var whisperModel = c.req.SpeechToText.Whisper.Model.String()
		if whisperModel != `` {
			var lang2 = c.req.AltLanguage
			var whisper = speech_to_text.NewWhisper(bibleId, c.database, whisperModel, lang2)
			status = whisper.ProcessFiles(audioFiles)
			if status.IsErr {
				return status
			}
			c.ident.TextSource = request.TextSTT
			if len(c.ident.AudioOTId) >= 10 {
				c.ident.TextOTId = c.ident.AudioOTId[:7] + `_TT`
			}
			if len(c.ident.AudioNTId) >= 10 {
				c.ident.TextNTId = c.ident.AudioNTId[:7] + `_TT`
			}
			c.database.UpdateIdent(c.ident)
		}
	}
	return status
}

func (c *Controller) encodeAudio(audioFiles []input.InputFile) dataset.Status {
	var status dataset.Status
	bibleId := c.req.BibleId
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

func (c *Controller) matchText() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	compare := match.NewCompare(c.ctx, c.user, c.req.Compare.BaseDataset, c.database, c.req.Testament, c.req.Compare.CompareSettings)
	filename, status = compare.Process()
	return filename, status
}

func (c *Controller) output() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	var out = output.NewOutput(c.ctx, c.database, c.req.DatasetName, false, false)
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
		}
	} else if c.req.OutputFormat.JSON {
		filename, status = out.WriteJSON(records, meta)
		if status.IsErr {
			return filename, status
		}
	}
	records = nil
	return filename, status
}

func (c *Controller) outputStatus(status dataset.Status) string {
	var filename string
	var status2 dataset.Status
	var out = output.NewOutput(c.ctx, db.DBAdapter{}, c.req.DatasetName, false, false)
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
