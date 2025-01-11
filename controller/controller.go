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
	"dataset/mms/fa_score_analysis"
	"dataset/output"
	"dataset/read"
	"dataset/request"
	"dataset/run_control"
	"dataset/speech_to_text"
	"dataset/timestamp"
	"time"
)

type OutputFiles struct {
	Directory string
	FilePaths []string
}

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
	log.Info(ctx, "Request: ", string(yamlContent))
	c.yamlRequest = yamlContent
	c.bucket = run_control.NewRunBucket(ctx, yamlContent)
	c.bucket.IsUnitTest = false // set to true when testing to make RunBucket work.
	return c
}

func (c *Controller) SetPostFiles(postFiles *input.PostFiles) {
	c.postFiles = postFiles
}

// Process is deprecated for production, but is a test only convenience method
func (c *Controller) Process() (string, dataset.Status) {
	output, status := c.ProcessV2()
	if status.IsErr {
		return "", status
	}
	if len(output.FilePaths) > 0 {
		return output.FilePaths[0], status
	} else {
		return "NO OUTPUT", status
	}
}

// ProcessV2 is the production means to execute the controller
func (c *Controller) ProcessV2() (OutputFiles, dataset.Status) {
	var start = time.Now()
	if c.postFiles != nil {
		defer c.postFiles.RemoveDir()
	}
	log.Debug(c.ctx)
	var status = c.processSteps()
	if status.IsErr {
		filename := c.outputStatus(status)
		c.bucket.AddOutput(filename)
	}
	var output OutputFiles
	output.Directory = c.req.Output.Directory
	output.FilePaths = c.bucket.GetOutputPaths()
	log.Info(c.ctx, "Duration", time.Since(start))
	log.Debug(c.ctx)
	c.bucket.PersistToBucket()
	return output, status
}

func (c *Controller) processSteps() dataset.Status {
	var filename string
	var status dataset.Status
	// Decode YAML Request File
	log.Info(c.ctx, "Parse .yaml file.")
	reqDecoder := request.NewRequestDecoder(c.ctx)
	c.req, status = reqDecoder.Process(c.yamlRequest)
	if status.IsErr {
		return status
	}
	c.ctx = context.WithValue(c.ctx, `request`, string(c.yamlRequest))
	// Get User
	log.Info(c.ctx, "Fetch Bible Brain data.")
	c.user, status = fetch.GetDBPUser(c.req)
	if status.IsErr {
		return status
	}
	// Open Database
	c.database, status = db.NewerDBAdapter(c.ctx, c.req.IsNew, c.user.Username, c.req.DatasetName)
	if status.IsErr {
		return status
	}
	defer c.database.Close()
	c.bucket.AddDatabase(c.database)
	// Fetch Ident Data from Ident
	c.ident, status = c.database.SelectIdent()
	if status.IsErr {
		return status
	}
	// Update Ident Data from DBP
	c.ident, status = c.fetchData()
	if status.IsErr {
		if c.req.TextData.AnyBibleBrain() || c.req.AudioData.AnyBibleBrain() {
			return status
		}
	}
	// Collect Text Input
	var textFiles []input.InputFile
	if !c.req.TextData.NoText {
		log.Info(c.ctx, "Load text files.")
		textFiles, status = c.collectTextInput()
		if status.IsErr {
			return status
		}
	}
	// Collect Audio Input
	var audioFiles []input.InputFile
	if !c.req.AudioData.NoAudio {
		log.Info(c.ctx, "Load audio files.")
		audioFiles, status = c.collectAudioInput()
		if status.IsErr {
			return status
		}
	}
	// Update Ident Table
	status = input.UpdateIdent(c.database, &c.ident, textFiles, audioFiles)
	if status.IsErr {
		return status
	}
	// Read Text Data
	if !c.req.TextData.NoText {
		log.Info(c.ctx, "Read and parse text files.")
		status = c.readText(textFiles)
		if status.IsErr {
			return status
		}
	}
	// Timestamps
	if !c.req.Timestamps.NoTimestamps {
		log.Info(c.ctx, "Read or create audio timestamp data.")
		status = c.timestamps(audioFiles)
		if status.IsErr {
			return status
		}
	}
	// Copy for STT
	if !c.req.TextData.NoText && !c.req.SpeechToText.NoSpeechToText {
		c.req.Compare.BaseDataset = c.database.Project
		c.req.AudioProof.BaseDataset = c.database.Project // ? should there be one BaseDataset ?
		// This makes a copy of database, and closes it.  Names the new database *_audio, and returns new
		c.database, status = c.database.CopyDatabase(`_audio`)
		if status.IsErr {
			return status
		}
		c.bucket.AddDatabase(c.database)
		status = c.database.UpdateEraseScriptText()
		if status.IsErr {
			return status
		}
	}
	// Speech to Text
	if !c.req.SpeechToText.NoSpeechToText {
		log.Info(c.ctx, "Perform speech to text.")
		status = c.speechToText(audioFiles)
		if status.IsErr {
			return status
		}
	}
	// Encode Audio
	if !c.req.AudioEncoding.NoEncoding {
		log.Info(c.ctx, "Perform audio encoding.")
		status = c.encodeAudio(audioFiles)
		if status.IsErr {
			return status
		}
	}
	// Encode Text
	if !c.req.TextEncoding.NoEncoding {
		log.Info(c.ctx, "Perform text encoding.")
		status = c.encodeText()
		if status.IsErr {
			return status
		}
	}
	// Audio Proofing
	if c.req.AudioProof.HTMLReport {
		log.Info(c.ctx, "Perform audio proof Report.")
		filename, status = c.audioProofing(audioFiles)
		if status.IsErr {
			return status
		}
		c.bucket.AddOutput(filename)
	}
	// Compare
	if c.req.Compare.HTMLReport {
		log.Info(c.ctx, "Perform text comparison.")
		filename, status = c.matchText()
		if status.IsErr {
			return status
		}
		c.bucket.AddOutput(filename)
	}
	// Prepare output
	log.Info(c.ctx, "Generate output.")
	if c.req.Output.Sqlite {
		c.bucket.AddOutput(c.database.DatabasePath)
	}
	if c.req.Output.CSV || c.req.Output.JSON {
		status = c.output()
		// added to bucket in c.output()
	}
	return status
}

func (c *Controller) fetchData() (db.Ident, dataset.Status) {
	var status dataset.Status
	var info fetch.BibleInfoType
	client := fetch.NewAPIDBPClient(c.ctx, c.req.BibleId)
	info, status = client.BibleInfo()
	if status.IsErr {
		c.ident, status = client.UpdateIdent(c.ident, info, c.req)
		return c.ident, status
	}
	client.FindFilesets(&info, c.req.AudioData.BibleBrain, c.req.TextData.BibleBrain, c.req.Testament)
	download := fetch.NewAPIDownloadClient(c.ctx, c.req.BibleId, c.req.Testament)
	status = download.Download(info)
	if status.IsErr {
		return c.ident, status
	}
	c.ident, status = client.UpdateIdent(c.ident, info, c.req)
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
		var ts mms.MMSAlign
		ts = mms.NewMMSAlign(c.ctx, c.database, c.ident.LanguageISO, c.req.AltLanguage)
		status = ts.ProcessFiles(audioFiles)
		if status.IsErr {
			return status
		}
		analysisRpt, _ := fa_score_analysis.FAScoreAnalysis(c.database)
		c.bucket.AddOutput(analysisRpt)
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

func (c *Controller) audioProofing(audioFiles []input.InputFile) (string, dataset.Status) {
	// Using audioFiles here should be tempoary, once the timestamps are updated with duration
	// there should be no need for the audio files to be present.
	var filename string
	var status dataset.Status
	if len(audioFiles) == 0 {
		return filename, log.ErrorNoErr(c.ctx, 400, "There are no audio files to AudioProof")
	}
	audioDir := audioFiles[0].Directory
	var textConn db.DBAdapter
	textConn, status = db.NewerDBAdapter(c.ctx, false, c.user.Username, c.req.AudioProof.BaseDataset)
	if status.IsErr {
		return filename, status
	}
	calc := match.NewAlignSilence(c.ctx, textConn, c.database) // c.database is ASR result
	faLines, filenameMap, status := calc.Process(audioDir)
	if status.IsErr {
		return filename, status
	}
	writer := match.NewAlignWriter(c.ctx, textConn)
	filename, status = writer.WriteReport(c.req.DatasetName, faLines, filenameMap)
	return filename, status
}

func (c *Controller) matchText() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	compare := match.NewCompare(c.ctx, c.user, c.req.Compare.BaseDataset, c.database, c.ident.LanguageISO, c.req.Testament, c.req.Compare.CompareSettings)
	filename, status = compare.Process()
	return filename, status
}

func (c *Controller) output() dataset.Status {
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
	if c.req.Output.CSV {
		filename, status = out.WriteCSV(records, meta)
		if status.IsErr {
			return status
		}
		c.bucket.AddOutput(filename)
	}
	if c.req.Output.JSON {
		filename, status = out.WriteJSON(records, meta)
		if status.IsErr {
			return status
		}
		c.bucket.AddOutput(filename)
	}
	records = nil
	return status
}

func (c *Controller) outputStatus(status dataset.Status) string {
	var filename string
	var status2 dataset.Status
	var out = output.NewOutput(c.ctx, db.DBAdapter{}, c.req.DatasetName, false, false)
	if c.req.Output.CSV {
		filename, status2 = out.CSVStatus(status, true)
	} else if c.req.Output.JSON {
		filename, status2 = out.JSONStatus(status, true)
	} else {
		filename = status.String()
	}
	if status2.IsErr {
		filename = status2.String()
	}
	return filename
}
