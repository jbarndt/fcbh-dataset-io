package tests

import (
	"testing"
)

// Test expects JMDYPM_audio.db and JMDYPM_text.db to exist.

const compareOnly = `is_new: no
dataset_name: JMDYPM_audio
bible_id: JMDYPM
username: GaryNTest
email: gary@shortsands.com
output_file: 15b_compare_only.html
testament:
  nt_books: [MAT,MRK,LUK,JHN,ACT]
compare:
  base_dataset: JMDYPM_text
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

func TestTwoCompareDirect(t *testing.T) {
	var tests []SqliteTest
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts", 26})
	//tests = append(tests, SqliteTest{"SELECT count(*) FROM scripts WHERE script_begin_ts != 0.0", 25})
	DirectSqlTest(compareOnly, tests, t)
}
