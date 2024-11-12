package tests

import (
	"strings"
	"testing"
)

const TSBibleBrain = `is_new: yes
dataset_name: TS_BB
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 13a_ts_bb.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps: 
  bible_brain: yes
testament:
  nt_books: ['MRK']
`

func TestTSBB(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 694})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 678})
	testName := strings.Replace(TSBibleBrain, "{bibleId}", "ENGWEB", -1)
	DirectSqlTest(testName, tests, t)
}

// ENGWEB BB timestamps
// select avg(script_end_ts-script_begin_ts) from scripts where script_end_ts != 0.0
// = 8.37511692230324
