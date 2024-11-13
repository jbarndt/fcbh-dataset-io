package tests

import (
	"strings"
	"testing"
)

const FastText = `is_new: yes
dataset_name: fasttext
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 12b_fasttext.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
text_encoding: 
  fast_text: yes
testament:
  nt: yes
detail:
  words: yes
`

func TestFasttextDirect(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 8218})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words WHERE ttype = 'W'", 175829})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM words WHERE word_enc != ''", 175829})
	testName := strings.Replace(FastText, "{bibleId}", "ENGWEB", -1)
	DirectSqlTest(testName, tests, t)
}
