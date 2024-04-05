package encode

import (
	"context"
	"dataset"
	"dataset/db"
	"testing"
)

func TestAeneas(t *testing.T) {
	var ctx = context.Background()
	var bibleId = `ENGWEB`
	var filesetId = `ENGWEBN2DA`
	var language = `eng`
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//files, status := ReadDirectory(ctx, bibleId, filesetId)
	aeneas := NewAeneas(ctx, conn, bibleId, filesetId)
	aeneas.Process(language, dataset.LINES)
}

/*
func TestAeneas(t *testing.T) {
	var bibleId = `ATIWBT`
	var filesetId = ``
	dir = "../../Desktop/Mark_Scott_1_1-31-2024/Audio Files"
	audioFile = "N2_MZI_BSM_046_LUK_002_VOX.wav"
	audioPath = os.path.join(dir, audioFile)
	textOutFile = "./aeneas_input.txt"
	if os.path.exists(textOutFile):
	os.remove(textOutFile)
	db = DBAdapter("ENG", 3, "Excel")
	createWordsFile(db, audioFile, textOutFile)
	outFile = "excel.json"
	aeneas("eng", audioPath, textOutFile, outFile)
	storeAeneas(db, audioFile, outFile)
}
*/
