package tests

import (
	"testing"
)

const proofCompare = `is_new: no
dataset_name: 15c_proof_compare_audio
bible_id: ENGWEB
username: GaryNTest
email: gary@shortsands.com
output_file: 15c_proof_compare.html
testament:
  nt: yes
audio_data:
  bible_brain:
    mp3_64: yes
audio_proof:
  html_report: yes
  base_dataset: 15c_proof_compare
compare:
  base_dataset: 15c_proof_compare
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

func TestProofAndCompare(t *testing.T) {
	CLIExec(proofCompare, t)
}

func TestProofAndCompareWithDebugger(t *testing.T) {
	DirectSqlTest(proofCompare, []SqliteTest{}, t)
}
