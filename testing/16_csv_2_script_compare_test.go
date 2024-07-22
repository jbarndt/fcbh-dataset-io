package testing

import (
	"dataset/request"
	"testing"
)

const LoadBaseScript = `is_new: yes
dataset_name: LoadBaseScript_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
output_file: 16__load_base_script.json
text_data:
  file: /Users/gary/FCBH2024/tugutil/Vessel Text_Tugutil_N2TUJNTM_03.xlsm
`

const CSV2ScriptCompare = `is_new: yes
dataset_name: CSV2ScriptCompare_{bibleId}
bible_id: {bibleId}
username: GaryNTest
email: gary@shortsands.com
#output_file: 16__csv_2_script_compare.html
output_file: 16__csv_2_script_compare_output.json
text_data:
  file: /Users/gary/FCBH2024/tugutil/TUJNTMN2TT.csv
testament:
  nt_books: [MRK]
`

func TestCSV2ScriptCompare(t *testing.T) {
	var tests []CtlTest
	//tests = append(tests, CtlTest{BibleId: "TUJNTM", Expected: 708, TextNtId: "TUJNTMN2ST",
	//	TextType: request.TextScript, AudioNTId: "", Language: "tvj"})
	//DirectTestUtility(LoadBaseScript, tests, t)
	//tests = nil
	tests = append(tests, CtlTest{BibleId: "TUJNTM", Expected: 788, TextNtId: "TUJNTMN2TT", // bibleId TUJNTM
		TextType: request.TextCSV, AudioNTId: "", Language: "tuj"})
	DirectTestUtility(CSV2ScriptCompare, tests, t)
}
