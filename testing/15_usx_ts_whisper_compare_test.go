package testing

import (
	"dataset/request"
	"testing"
)

const USXTSWhisperCompare = `is_new: yes
dataset_name: USXTSWhisperCompare_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 15__usx_ts_whisper_compare.html
text_data:
  bible_brain:
    text_usx_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  aeneas: yes
testament:
  nt_books: ['TIT']
speech_to_text:
  language: en
  whisper:
    model:
      tiny: yes
compare:
  compare_settings: 
    lower_case: y
    remove_prompt_chars: y
    remove_punctuation: y
    double_quotes: 
      remove: y
    apostrophe: 
      remove: y
    hyphen:
      remove: y
    diacritical_marks:
      normalize_nfd: y
`

func TestUSXTSWhisperCompare(t *testing.T) {
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 27, TextNtId: "ENGWEBN_ET-usx",
		TextType: request.TextUSXEdit, AudioNTId: "ENGWEBN2DA-mp3-16", Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(USXTSWhisperCompare, tests, t)
}

func TestPlainWhisperCompare(t *testing.T) {
	const PlainTSWhisperCompare = `is_new: yes
dataset_name: PlainTSWhisperCompare_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 15__plain_ts_whisper_compare.html
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  aeneas: yes
testament:
  nt_books: ['TIT']
compare:
  base_dataset: USXTSWhisperCompare_{bibleId}_STT
  compare_settings: 
    lower_case: y
    remove_prompt_chars: y
    remove_punctuation: y
    double_quotes: 
      remove: y
    apostrophe: 
      remove: y
    hyphen:
      remove: y
    diacritical_marks:
      normalize_nfd: y
`
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 27, TextNtId: "ENGWEBN_ET",
		TextType: request.TextPlainEdit, Language: "eng"})
	//tests = append(tests, try{bibleId: "ATIWBT", expected: 7, textNtId: "ATIWBTN_ET", audioNTId: "ATIWBTN1DA",
	//	language: "ati"}) // There are no timestamps
	DirectTestUtility(PlainTSWhisperCompare, tests, t)
}
