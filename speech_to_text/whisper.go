package speech_to_text

import (
	"bytes"
	"dataset_io"
	"dataset_io/db"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Docs:
https://github.com/openai/whisper
Install:
pip3 install git+https://github.com/openai/whisper.git
Whisper is an open source Speech to Text program developed by OpenAI.
Executable:
/Users/gary/Library/Python/3.9/bin/whisper

1. Retry Whisper with a non-drama version of the audio.
2. Pi recommends preprocessing in audacity
3. Do not get text, but get segments, and iterate over it.
4. Capture start, end timestamps
5. Capture text
6. Capture tokens, if I have a place for it
7. Capture avg_logprob
8. Capture no_speech_prob
9. Capture compression ratio

'segments': [{'id': 0, 'seek': 0, 'start': 0.0, 'end': 3.24, 'text': ' Chapter 3', 'tokens': [50363, 7006, 513, 50525],
'temperature': 0.0, 'avg_logprob': -0.2316610102067914, 'compression_ratio': 1.46, 'no_speech_prob':
0.2119932472705841},

*/

type Whisper struct {
	conn    db.DBAdapter
	records []db.InsertScriptRec
}

func NewWhisper(conn db.DBAdapter) Whisper {
	var w Whisper
	w.conn = conn
	w.records = make([]db.InsertScriptRec, 0, 100000)
	return w
}

func (w *Whisper) ProcessDirectory(bibleId string, filesetId string, testament dataset_io.TestamentType) {
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, filesetId)
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files { // in python sorted(os.listdir(directory))
		filename := file.Name()
		if !strings.HasPrefix(filename, `.`) {
			fmt.Println(filename)
			fileType := filename[:1]
			if fileType == `A` && (testament == dataset_io.OT || testament == dataset_io.ONT) {
				w.processFile(directory, filename)
			} else if fileType == `B` && (testament == dataset_io.NT || testament == dataset_io.ONT) {
				w.processFile(directory, filename)
			}
		}
	}
}

func (w *Whisper) processFile(directory string, filename string) {
	bookId, chapter := w.parseFilename(filename)
	if bookId == `TIT` && chapter == 3 {
		var path = filepath.Join(directory, filename)
		w.runWhisper(path)
	}

	//result = self.model.transcribe(path)
	//scriptText := result["text"]
	//var rec db.InsertScriptRec
	//bookId, chapter := w.parseFilename(filename)
	//rec.BookId = bookId
	//rec.ChapterNum = chapter
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//rec.AudioFile = filename
	//rec.VerseNum = 1 // ????
	//rec.VerseStr = `1`
	//rec.ScriptText = scriptText
	//w.records = append(w.records, rec)
}

func (w *Whisper) parseFilename(filename string) (string, int) {
	chapter, err := strconv.Atoi(filename[6:8])
	if err != nil {
		log.Fatal(err)
	}
	bookName := strings.Replace(filename[9:21], `_`, ``, -1)
	bookId := db.USFMBookId(bookName)
	return bookId, chapter
}

/*
/Users/gary/Library/Python/3.9/bin/whisper
--model medium //large small
--model_dir Pathname defaults to .
--output_dir  default .
--output_format txt,vtt,tsv,json,all all is default
--task transcript || translate
--language
*/
func (w *Whisper) runWhisper(audioFilePath string) {
	whisperPath := `/Users/gary/Library/Python/3.9/bin/whisper`
	//model := `--model ` + `tiny`
	//format := `--output_format json`
	cmd := exec.Command(whisperPath, audioFilePath, `--model`, `tiny`, `--output_format`, `json`)
	fmt.Println(cmd.String())
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	stderrStr := stderrBuf.String()
	if stderrStr != `` {
		fmt.Printf("Stderr: \n%s\n", stderrStr)
	}
	fmt.Println("\n\nSTDOUT:", stdoutBuf.String())
	//if stdoutStr != "" {
	//	fmt.Printf("Stdout: \n%s\n", stdoutStr)
	//}
	os.Exit(0)

	type WhisperSegmentType struct {
		Id     int     `json:"id"`
		Seek   float32 `json:"seek"`
		Start  float32 `json:"start"`
		End    float32 `json:"end"`
		Text   string  `json:"text"`
		Tokens []int   `json:"tokens"`
		// "tokens": [50660, 293, 281, 12076, 11, 281, 312, 42541, 11, 281, 312, 1919, 337, 633, 665, 589, 11, 281, 1710, 6724, 295, 51004],
		Temperature      float32 `json:"temperature"`
		AvgLogProb       float64 `json:"avg_logprob"`
		CompressionRatio float64 `json:"compression_ratio"`
		NoSpeechProb     float64 `json:"no_speech_prob"`
	}
	type WhisperOutputType struct {
		Segments []WhisperSegmentType `json:"segments"`
		Language string               `json:"language"`
	}
	var response WhisperOutputType
	err = json.Unmarshal(stdoutBuf.Bytes(), &response)
	if err != nil {
		log.Fatalln("Error decoding Whisper JSON:", err)
	}
	for _, seg := range response.Segments {
		fmt.Println(seg)
	}
}
