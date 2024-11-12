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

/*
func TestMFCCLines(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var testament = request.Testament{NTBooks: []string{`MRK`}}
	testament.BuildBookMaps()
	var detail = request.Detail{Lines: true}
	files, status := input.DBPDirectory(ctx, bibleId, `audio`, ``, filesetId, testament)
	if status.IsErr {
		t.Error(status.Message)
	}
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	conn.DeleteMFCCs()
	mfcc := NewMFCC(ctx, conn, bibleId, detail, 7)
	status = mfcc.ProcessFiles(files)
	if status.IsErr {
		t.Error(status.Message)
	}
	count, _ := conn.CountScriptMFCCRows()
	if count != 678 {
		t.Error(`Script count should be 678`, count)
	}
	count, _ = conn.CountWordMFCCRows()
	if count != 0 {
		t.Error(`Word count should be 0`, count)
	}
}

*/
