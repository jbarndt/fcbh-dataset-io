package req

import (
	"fmt"
	"testing"
)

func TestRequestYamlFile(t *testing.T) {
	req, status := DecodeFile(`request.yaml`)
	fmt.Println(`Status:`, status)
	yaml, status := Encode(req)
	fmt.Println(yaml)
}

func TestParser(t *testing.T) {
	var test1 = `Required:
  RequestName: Test1  # should be a unique name
  RequestorName: Gary G
  RequestorEmail: gary@shortsands.com
  BibleId: ENGWEB
  LanguageISO: eng
  VersionCode: WEB
Testament:
  OT: Y
  NT: Yes
AudioData:
  BibleBrain:
    MP3_64: Yes
    MP3_16: Y
    OPUS: y
  File: file:///Users/who/where.mp3
  Http: http://go.there.com/path
  AWSS3: http://west1/path
  POST: Yes
  NoAudio: Yes
TextData: 
  BibleBrain: 
    TextUSXEdit: Yes
    TextPlainEdit: Yes
    TextPlain: Yes
  SpeechToText: 
    Whisper: 
      Model: 
        Large: Yes
        Medium: Yes
        Small: Yes
        Base: Yes
        Tiny: Yes
  File: /Users/who/where.mp3
  Http: http://go.there.com/path
  AWSS3: http://west1/path
  POST: Yes
  NoText: Yes
Detail:
  Lines: yes
  Words: yes
Timestamps:
  BibleBrain: yes
  Aeneas: yes
  NoTimestamps: yes
AudioEncoding:
  MFCC: yes
  NoEncoding: yes
TextEncoding:
  FastText: yes
  NoEncoding: yes
OutputFormat:
  CSV: yes
  JSON: yes
  Sqlite: yes
Compare:
  Project1: UseProject1
  Project2: UseProject2
  CompareSettings:
    LowerCase: yes
    RemovePromptChars: yes
    RemovePunctuation: yes
    DoubleQuotes: 
      Remove: yes
      Normalize: yes
    Apostrophe: 
      Remove: yes
      Normalize: yes
    Hyphen: 
      Remove: yes
      Normalize: yes
    DiacriticalMarks: 
      Remove: yes
      NormalizeNFC: yes
      NormalizeNFD: yes
      NormalizeNFKC: yes
      NormalizeNFKD: yes`
	req, status := DecodeString(test1)
	fmt.Println(`Status:`, status)
	if !req.TextEncoding.FastText {
		t.Error("FastText should be true")
	}
	_, _ = Encode(req)
	var boolTests = []bool{req.Testament.OT, req.Testament.NT,
		req.AudioData.BibleBrain.MP3_64,
		req.AudioData.BibleBrain.MP3_16,
		req.AudioData.BibleBrain.OPUS,
		req.AudioData.POST,
		req.AudioData.NoAudio,
		req.TextData.BibleBrain.TextUSXEdit,
		req.TextData.BibleBrain.TextPlainEdit,
		req.TextData.BibleBrain.TextPlain,
		req.TextData.SpeechToText.Whisper.Model.Large,
		req.TextData.SpeechToText.Whisper.Model.Medium,
		req.TextData.SpeechToText.Whisper.Model.Small,
		req.TextData.SpeechToText.Whisper.Model.Base,
		req.TextData.SpeechToText.Whisper.Model.Tiny,
		req.TextData.POST,
		req.TextData.NoText,
		req.Detail.Lines,
		req.Detail.Words,
		req.Timestamps.BibleBrain,
		req.Timestamps.Aeneas,
		req.Timestamps.NoTimestamps,
		req.AudioEncoding.MFCC,
		req.AudioEncoding.NoEncoding,
		req.TextEncoding.FastText,
		req.TextEncoding.NoEncoding,
		req.OutputFormat.CSV,
		req.OutputFormat.JSON,
		req.OutputFormat.Sqlite,
		req.Compare.CompareSettings.LowerCase,
		req.Compare.CompareSettings.RemovePromptChars,
		req.Compare.CompareSettings.RemovePunctuation,
		req.Compare.CompareSettings.DoubleQuotes.Remove,
		req.Compare.CompareSettings.DoubleQuotes.Normalize,
		req.Compare.CompareSettings.Apostrophe.Remove,
		req.Compare.CompareSettings.Apostrophe.Normalize,
		req.Compare.CompareSettings.Hyphen.Remove,
		req.Compare.CompareSettings.Hyphen.Normalize,
		req.Compare.CompareSettings.DiacriticalMarks.Remove,
		req.Compare.CompareSettings.DiacriticalMarks.NormalizeNFC,
		req.Compare.CompareSettings.DiacriticalMarks.NormalizeNFD,
		req.Compare.CompareSettings.DiacriticalMarks.NormalizeNFKC,
		req.Compare.CompareSettings.DiacriticalMarks.NormalizeNFKD}
	for i, item := range boolTests {
		if !item {
			t.Error(`The`, i, `th item should be true, but is not`)
		}
	}
	var strs = []string{req.Required.RequestName,
		req.Required.RequestorName,
		req.Required.RequestorEmail,
		req.Required.BibleId,
		req.Required.LanguageISO,
		req.Required.VersionCode,
		req.AudioData.File,
		req.AudioData.Http,
		req.AudioData.AWSS3,
		req.TextData.File,
		req.TextData.Http,
		req.TextData.AWSS3,
		req.Compare.Project1,
		req.Compare.Project2}
	for i, item := range strs {
		if len(item) == 0 {
			t.Error(`The`, i, `th item should have a value, but is empty`)
		}
	}
}
