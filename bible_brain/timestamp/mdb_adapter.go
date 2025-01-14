package timestamp

import (
	"database/sql"
	"dataset"
	"fmt"
)

// https://github.com/faithcomesbyhearing/verse-timing/blob/master/BibleFileTimestamps_Insert_aeneas.py

// include here the code to update mdb, including opening the connection and closing
// the connection.

func OpenMDB() *sql.DB {
	var conn *sql.DB
	return conn
}

func SelectBibleIds() {

}

func SelectFilesets(bibleId string) {

}

func SelectFileId(filesetId string, bookId string, chapterNum int) (int, dataset.Status) {
	var result int
	var status dataset.Status
	query := `SELECT distinct bf.id FROM bible_files bf, bible_filesets bs
			WHERE bf.hash_id =  bs.hash_id
			AND bs.id = ? 
			AND bf.book_id = ?
			AND bf.chapter_start = ?`
	fmt.Println(query, filesetId, bookId, chapterNum)
	return result, status
}

func DeleteTimestamps(filesetId string) dataset.Status {
	var status dataset.Status
	query := `DELETE FROM bible_file_timestamps WHERE bible_file_id IN
		(SELECT bf.id FROM bible_files bf, bible_filesets bs"
		WHERE bf.hash_id = bs.hash_id AND bs.id = ?)`
	fmt.Println(query, filesetId)
	return status
}

func SelectHashId(filesetId string) (string, dataset.Status) {
	var result string
	var status dataset.Status
	query := `SELECT hash_id FROM bible_filesets WHERE asset_id = 'dbp-prod'
		AND set_type_code IN ('audio', 'audio_drama') AND id = ?")`
	fmt.Println(query, filesetId)
	return result, status
}

func UpdateFilesetTimingEstTag(hashId string, timingEstErr int) dataset.Status {
	var status dataset.Status
	// probably should use 10 to mean fa_align
	query := `REPLACE INTO bible_fileset_tags (hash_id, name, description, admin_only, iso, language_id)
		VALUES (?, 'timing_est_err', ?, 0, 'eng', 6414)")`
	fmt.Println(query, hashId, timingEstErr)
	return status
}

type Timestamp struct {
	BibleFileId int
	VerseStr    string
	VerseEnd    string
	BeginTS     float64
	EndTS       float64
}

func UpdateTimestamps(timestamps []Timestamp) dataset.Status {
	var status dataset.Status
	query := `INSERT INTO bible_file_timestamps (bible_file_id, verse_start, verse_end,
		timestamp) VALUES (?, ?, ?, ?)`
	fmt.Println(query)
	return status
}
