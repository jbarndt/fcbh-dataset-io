package req

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"testing"
)

func TestRequestYamlFile(t *testing.T) {
	req := DecodeFile(`request.yaml`)
	encode(req)
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
	req := DecodeString(test1)
	if !req.TextEncoding.FastText {
		t.Error("FastText should be true")
	}
	encode(req)
}

func encode(req Request) {
	d, err := yaml.Marshal(&req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}
