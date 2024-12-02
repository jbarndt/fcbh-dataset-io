package tests

import (
	"testing"
)

const mmsAlignTest = `is_new: yes
dataset_name: 13e_mms_align
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: 13e_mms_align.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  mms_align: yes
testament:
  nt_books: [PHM]
`

func TestMMSAlignDirect(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words", 448})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words WHERE word_begin_ts != 0.0", 448})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM chars", 1796})
	tests = append(tests, SqliteTest{"SELECT count(distinct(word_id)) FROM chars", 448})
	DirectSqlTest(mmsAlignTest, tests, t)
}
