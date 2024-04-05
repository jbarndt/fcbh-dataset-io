package encode

import (
	"bytes"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Aeneas struct {
	ctx       context.Context
	conn      db.DBAdapter
	bibleId   string
	audioFSId string
}

// NewAeneas the DBAdapter contains text to be timestamped.
func NewAeneas(ctx context.Context, conn db.DBAdapter, bibleId string, audioFSId string) Aeneas {
	var a Aeneas
	a.ctx = ctx
	a.conn = conn
	a.bibleId = bibleId
	a.audioFSId = audioFSId
	return a
}

func (a *Aeneas) Process(language string, detail dataset.TextDetailType) dataset.Status {
	var status dataset.Status
	audioFiles, status := ReadDirectory(a.ctx, a.bibleId, a.audioFSId)
	if status.IsErr {
		return status
	}
	if detail == dataset.LINES || detail == dataset.BOTH {
		status = a.processScripts(language, audioFiles)
	} else if detail == dataset.WORDS || detail == dataset.BOTH {
		status = a.processWords(language, audioFiles)
	}
	return status
}

func (a *Aeneas) processScripts(language string, audioFiles []string) dataset.Status {
	var status dataset.Status
	for _, audioFile := range audioFiles {
		bookId, chapterNum, status := ParseFilename(a.ctx, audioFile)
		if status.IsErr {
			return status
		}
		scripts, status := a.conn.SelectScriptsByBookChapter(bookId, chapterNum)
		if status.IsErr {
			return status
		}
		var scriptText = make([]string, 0, len(scripts))
		for _, script := range scripts {
			scriptText = append(scriptText, script.ScriptText)
		}
		textFile, status := a.createFile(bookId, chapterNum, scriptText)
		if status.IsErr {
			return status
		}
		fmt.Println(textFile, status)
		outputFile, status := a.executeAeneas(language, audioFile, textFile)
		if status.IsErr {
			return status
		}
		fmt.Println("Output", outputFile)
		fragments, status := a.parseResponse(outputFile)
		for _, frag := range fragments {
			fmt.Println(frag)
		}
		scripts, status = a.mergeScriptTimestamps(audioFile, scripts, fragments)
		if status.IsErr {
			return status
		}
		status = a.conn.UpdateScriptTimestamps(scripts)
		if status.IsErr {
			return status
		}
		scripts = nil
	}
	return status
}

func (a *Aeneas) processWords(language string, audioFiles []string) dataset.Status {
	var status dataset.Status
	for _, audioFile := range audioFiles {
		bookId, chapterNum, status := ParseFilename(a.ctx, audioFile)
		if status.IsErr {
			return status
		}
		words, status := a.conn.SelectWordsByBookChapter(bookId, chapterNum)
		if status.IsErr {
			return status
		}
		var wordText = make([]string, 0, len(words))
		for _, word := range words {
			wordText = append(wordText, word.Word)
		}
		textFile, status := a.createFile(bookId, chapterNum, wordText)
		if status.IsErr {
			return status
		}
		fmt.Println(textFile, status)
		outputFile, status := a.executeAeneas(language, audioFile, textFile)
		if status.IsErr {
			return status
		}
		fmt.Println("Output", outputFile)
		fragments, status := a.parseResponse(outputFile)
		for _, frag := range fragments {
			fmt.Println(frag)
		}
		words, status = a.mergeWordTimestamps(words, fragments)
		if status.IsErr {
			return status
		}
		status = a.conn.UpdateWordTimestamps(words)
		if status.IsErr {
			return status
		}
		words = nil
	}
	return status
}

func (a *Aeneas) createFile(bookId string, chapter int, scripts []string) (string, dataset.Status) {
	var filename string
	var status dataset.Status
	var fp, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), a.audioFSId+bookId+strconv.Itoa(chapter)+`_`)
	fmt.Println("Created file", fp.Name())
	if err != nil {
		status = log.Error(a.ctx, 500, err, `Unable to open temp file for scripts`)
		return filename, status
	}
	for _, script := range scripts {
		text := strings.Replace(script, "\n", ` `, -1)
		text = strings.TrimSpace(text)
		fp.WriteString(text)
		fp.WriteString("\n")
	}
	fp.Close()
	return fp.Name(), status
}

func (a *Aeneas) executeAeneas(language string, audioFile string, textFile string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	fname := filepath.Base(audioFile)
	fname = strings.Split(fname, `.`)[0]
	var output, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), fname+`_`)
	if err != nil {
		status = log.Error(a.ctx, 500, err, `Error creating temp output file in Aeneas`)
		return result, status
	}
	//pythonPath := os.Getenv(`PYTHON_EXE`)
	pythonPath := "python3"
	cmd := exec.Command(pythonPath, `-m`, `aeneas.tools.execute_task`,
		audioFile,
		textFile,
		`task_language=`+language+`|os_task_file_format=json|is_text_type=plain`,
		output.Name(),
		`-example-words --presets-word`)
	fmt.Println(cmd.String())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	if err != nil {
		status = log.Error(a.ctx, 500, err, `Error executing Aeneas`)
		return output.Name(), status
	}
	if stderrBuf.Len() > 0 {
		fmt.Println("STDERR", stderrBuf.String())
	}
	if stdoutBuf.Len() > 0 {
		fmt.Println("STDOUT", stdoutBuf.String())
	}
	return output.Name(), status
}

type AeneasRec struct {
	Begin    string   `json:"begin"`
	Children []string `json:"children"`
	End      string   `json:"end"`
	Id       string   `json:"id"`
	Language string   `json:"language"`
	Lines    []string `json:"lines"`
}

type AeneasResp struct {
	Fragments []AeneasRec `json:"fragments"`
}

func (a *Aeneas) parseResponse(filename string) ([]AeneasRec, dataset.Status) {
	var results []AeneasRec
	var status dataset.Status
	var content, err = os.ReadFile(filename)
	if err != nil {
		status = log.Error(a.ctx, 500, err, `Failed to open Aeneas output file`)
		return results, status
	}
	var response AeneasResp
	err = json.Unmarshal(content, &response)
	if err != nil {
		status = log.Error(a.ctx, 500, err, `Error parsing Aeneas output json`)
		return results, status
	}
	return response.Fragments, status
}

func (a *Aeneas) mergeScriptTimestamps(audioFile string, scripts []db.Script, fragments []AeneasRec) ([]db.Script, dataset.Status) {
	var status dataset.Status
	if len(scripts) != len(fragments) {
		status = log.ErrorNoErr(a.ctx, 500, `Scripts len`, len(scripts), `Aeneas len`, len(fragments))
		return scripts, status
	}
	var err error
	for i, scp := range scripts {
		frag := fragments[i]
		scp.AudioFile = filepath.Base(audioFile)
		scp.ScriptBeginTS, err = strconv.ParseFloat(frag.Begin, 64)
		if err != nil {
			status = log.Error(a.ctx, 500, err, `Could not parse begin TS from Aeneas`)
			return scripts, status
		}
		scp.ScriptEndTS, err = strconv.ParseFloat(frag.End, 64)
		if err != nil {
			status = log.Error(a.ctx, 500, err, `Could not parse end TS from Aeneas`)
			return scripts, status
		}
		scripts[i] = scp
	}
	return scripts, status
}

func (a *Aeneas) mergeWordTimestamps(words []db.Word, fragments []AeneasRec) ([]db.Word, dataset.Status) {
	var status dataset.Status
	if len(words) != len(fragments) {
		status = log.ErrorNoErr(a.ctx, 500, `Scripts len`, len(words), `Aeneas len`, len(fragments))
		return words, status
	}
	var err error
	for i, word := range words {
		frag := fragments[i]
		word.WordBeginTS, err = strconv.ParseFloat(frag.Begin, 64)
		if err != nil {
			status = log.Error(a.ctx, 500, err, `Could not parse begin TS from Aeneas`)
			return words, status
		}
		word.WordEndTS, err = strconv.ParseFloat(frag.End, 64)
		if err != nil {
			status = log.Error(a.ctx, 500, err, `Could not parse end TS from Aeneas`)
			return words, status
		}
		words[i] = word
	}
	return words, status
}
