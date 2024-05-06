package request

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestRequestYamlFile(t *testing.T) {
	var d = NewRequestDecoder(context.Background())
	content, err := os.ReadFile(`../controller/client/request_test.yaml`)
	if err != nil {
		panic(err)
	}
	req, status := d.Decode(content)
	fmt.Println(`Status:`, status)
	yaml, status := d.Encode(req)
	fmt.Println(yaml)
}

func TestParser(t *testing.T) {
	var test1 = `is_new: Yes
dataset_name: Test1  # should be a unique name
bible_id: ENGWEB
testament:
  ot: 
  ot_books: [GEN,EXO,LEV,NUM,DEU]
  nt: Yes
  nt_books: []
audio_data:
  bible_brain:
    mp3_64: Yes
    mp3_16: Y
    opus: y
  file: file:///Users/who/where.mp3
  aws_s3: http://west1/path
  post: Yes
  no_audio: Yes
text_data: 
  bible_brain: 
    text_usx_edit: Yes
    text_plain_edit: Yes
    text_plain: Yes
  speech_to_text: 
    whisper: 
      model: 
        large: Yes
        medium: Yes
        small: Yes
        base: Yes
        tiny: Yes
  file: /Users/who/where.mp3
  aws_s3: http://west1/path
  post: Yes
  no_text: Yes
detail:
  lines: yes
  words: yes
timestamps:
  bible_brain: yes
  aeneas: yes
  no_timestamps: yes
audio_encoding:
  mfcc: yes
  no_encoding: yes
text_encoding:
  fast_text: yes
  no_encoding: yes
output_format:
  csv: yes
  json: yes
  sqlite: yes
compare:
  base_dataset: UseProject1
  compare_settings:
    lower_case: yes
    remove_prompt_chars: yes
    remove_punctuation: yes
    double_quotes: 
      remove: yes
      normalize: yes
    apostrophe: 
      remove: yes
      normalize: yes
    hyphen: 
      remove: yes
      normalize: yes
    diacritical_marks: 
      remove: yes
      normalize_nfc: yes
      normalize_nfd: yes
      normalize_nfkc: yes
      normalize_nfkd: yes`
	var r = NewRequestDecoder(context.Background())
	req, status := r.Decode([]byte(test1))
	fmt.Println(`Status:`, status)
	if !req.TextEncoding.FastText {
		t.Error("FastText should be true")
	}
	_, _ = r.Encode(req)
	var boolTests = []bool{
		//req.Testament.OT,
		req.Testament.NT,
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
	var strs = []string{req.RequestName,
		req.BibleId,
		req.AudioData.File,
		req.AudioData.AWSS3,
		req.TextData.File,
		req.TextData.AWSS3,
		req.Compare.BaseProject}
	for i, item := range strs {
		if len(item) == 0 {
			t.Error(`The`, i, `th item should have a value, but is empty`)
		}
	}
	if len(req.Testament.OTBooks) != 5 {
		t.Error(`OTBooks should have a length of 5, not`, len(req.Testament.OTBooks))
	}
	if len(req.Testament.NTBooks) != 0 {
		t.Error(`NTBooks should have a length of 0, not`, len(req.Testament.NTBooks))
	}
}
