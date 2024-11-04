package testing

import (
	"testing"
)

const TSBucketTest = `is_new: yes
dataset_name: 13c_ts_bucket
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: 13c_ts_bucket.sqlite
text_data:
  bible_brain:
    text_plain_edit: yes
audio_data:
  bible_brain:
    mp3_64: yes
timestamps:
  ts_bucket: yes
testament:
  nt_books: [TIT]
`

func TestTSBucketDirect(t *testing.T) {
	var tests []SqliteTest
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 49})
	tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 46})
	DirectSqlTest(TSBucketTest, tests, t)
}
