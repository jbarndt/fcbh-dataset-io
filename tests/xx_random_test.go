package tests

import (
	"testing"
)

const test1 = `is_new: yes
dataset_name: ENGWEB_align_wav
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: ENGWEB_align_wav.json
text_data:
  bible_brain:
    text_usx_edit: yes
audio_data:
  aws_s3: s3://pretest-audio/N2ENGWEB English (ENG)/N2ENGWEB Chapter VOX/*.wav
timestamps:
  mms_fa_verse: yes
speech_to_text:
  mms_asr: yes
testament:
  nt: 
  nt_books: [LUK]
`

func TestRandomDirect(t *testing.T) {
	var tests []SqliteTest
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	DirectSqlTest(test1, tests, t)
}
