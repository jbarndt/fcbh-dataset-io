package tests

import (
	"dataset/request"
	"testing"
)

const MMSASRCompare = `is_new: yes
dataset_name: MMSASRCompare
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: 15a_mms_asr_compare.html
text_data:
  bible_brain:
    text_usx_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  bible_brain: yes
testament:
  nt_books: ['3JN']
speech_to_text:
  mms_asr: yes
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

func TestMMSASRCompare(t *testing.T) {
	var tests []CtlTest
	tests = append(tests, CtlTest{BibleId: "ENGWEB", Expected: 8, TextNtId: "ENGWEBN_ET-usx",
		TextType: request.TextUSXEdit, AudioNTId: "ENGWEBN2DA", Language: "eng"})
	//tests = append(tests, CtlTest{BibleId: "APFCMU", Expected: 16, TextNtId: "APFCMUN_ET-usx",
	//	AudioNTId: `APFCMUN1DA`, TextType: request.TextUSXEdit, Language: "apf"})
	DirectTestUtility(MMSASRCompare, tests, t)
}
