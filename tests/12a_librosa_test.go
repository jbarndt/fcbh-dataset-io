package tests

import (
	"strings"
	"testing"
)

const Librosa = `is_new: yes
dataset_name: librosa
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 12a_librosa.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps: 
  bible_brain: yes
audio_encoding: 
  mfcc: yes
testament:
  nt_books: [MRK]
`

func TestLibrosaDirect(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 694})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 678})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM script_mfcc", 694})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM word_mfcc", 0})
	testName := strings.Replace(Librosa, "{bibleId}", "ENGWEB", -1)
	DirectSqlTest(testName, tests, t)
}
