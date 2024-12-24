package tests

import (
	"testing"
)

const test1 = `is_new: no
dataset_name: N2ENG_WEB_WAV_audio
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: ENGWEB_WAV_vs_MP3.html
testament:
  nt: 
  nt_books: [LUK]
compare: 
  base_dataset: N2ENG_WEB_MP3_audio
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

const test2 = `is_new: yes
dataset_name: N2KTB_ESB
bible_id: KTBESB
username: GaryNTest
email: gary@shortsands.com
output_file: N2KTB_ESB.html
#alt_language: ktb
text_data:
  aws_s3: s3://pretest-audio/N2KTBESB Kambaata (KTB)/N2KTBESB Text/USX/*.usx
audio_data:
  aws_s3: s3://pretest-audio/N2KTBESB Kambaata (KTB)/N2KTBESB Chapter VOX/MP3/*.mp3
timestamps:
  mms_fa_verse: yes
testament:
  nt: 
  nt_books: [MRK]
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
const test3 = `is_new: yes
dataset_name: N2ENGWEB
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: N2ENGWEB.sqlite
text_data:
  bible_brain:
    text_usx_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  mms_align: yes
testament:
  nt: yes
speech_to_text:
  mms_asr: yes
`

func TestRandomDirect(t *testing.T) {
	var tests []SqliteTest
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	DirectSqlTest(test3, tests, t)
}
