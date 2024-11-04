package testing

import (
	"testing"
)

const FAWordTest = `is_new: yes
dataset_name: 13e_fa_word
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: 13e_fa_word.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  mms_fa_word: yes
testament:
  nt_books: [PHM]
`

func TestFAVWordDirect(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words", 447})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words WHERE word_begin_ts != 0.0", 447})
	DirectSqlTest(FAWordTest, tests, t)
}
