package tests

import (
	"testing"
)

const test1 = `is_new: yes
dataset_name: O2LISNVS
bible_id: NVSLIS
username: GaryNTest
email: gary@shortsands.com
output_file: O2LISNVS.json
text_data:
  aws_s3: s3://pretest-audio/LISNVS [T]/O2LISNVS Transliterated Text/*.usx
audio_data:
  aws_s3: s3://pretest-audio/LISNVS [T]/LISNVS Chapter VOX/*.mp3
#timestamps:
#  mms_fa_verse: yes
testament:
  ot: yes
`

func TestRandomDirect(t *testing.T) {
	var tests []SqliteTest
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	DirectSqlTest(test1, tests, t)
}
