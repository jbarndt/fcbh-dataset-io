package output

import (
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
)

func LoadScriptStruct(d db.DBAdapter) []Script {
	var results []Script
	var status dataset.Status
	query := `SELECT scripts.script_id, book_id, chapter_num, chapter_end, audio_file, script_num, 
		usfm_style, person, actor, verse_str, verse_end, script_text, 
		script_begin_ts, script_end_ts, rows, cols, mfcc_json
		FROM scripts LEFT OUTER JOIN script_mfcc ON script_mfcc.script_id = scripts.script_id
		WHERE book_id = 'MRK' AND chapter_num = 1 ORDER BY scripts.script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during select scripts")
		panic(status.Message)
	}
	defer rows.Close()
	for rows.Next() {
		var sc Script
		var mfccJson string
		err := rows.Scan(&sc.ScriptId, &sc.BookId, &sc.ChapterNum, &sc.ChapterEnd, &sc.AudioFile,
			&sc.ScriptNum, &sc.UsfmStyle, &sc.Person, &sc.Actor, &sc.VerseStr, &sc.VerseEnd,
			&sc.ScriptText, &sc.ScriptBeginTS, &sc.ScriptEndTS, &sc.MFCCRows, &sc.MFCCCols, &mfccJson)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error in SelectScripts.")
			//return results, status
			panic(status.Message)
		}
		err = json.Unmarshal([]byte(mfccJson), &sc.MFCC)
		if err != nil {
			panic(err)
		}
		results = append(results, sc)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results
}
