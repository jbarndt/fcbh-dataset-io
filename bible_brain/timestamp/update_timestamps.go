package timestamp

import (
	"context"
	"dataset"
	"dataset/bible_brain"
	"dataset/db"
)

const (
	mmsAlignTimingEstErr = "mms_align"
)

type UpdateTimestamps struct {
	ctx     context.Context
	conn    db.DBAdapter
	dbpConn bible_brain.DBPAdapter
}

func (d *UpdateTimestamps) Process(filesetId string) dataset.Status {
	var status dataset.Status
	status = d.dbpConn.DeleteTimestamps(filesetId)
	if status.IsErr {
		return status
	}
	var hashId string
	hashId, status = d.dbpConn.SelectHashId(filesetId)
	if status.IsErr {
		return status
	}
	status = d.dbpConn.UpdateFilesetTimingEstTag(hashId, mmsAlignTimingEstErr)
	if status.IsErr {
		return status
	}
	var books []string
	books = append(books, db.BookOT...)
	books = append(books, db.BookNT...)
	for _, book := range books {
		lastChapter, _ := db.BookChapterMap[book]
		for chap := 1; chap <= lastChapter; chap++ {
			var timestamps []db.Audio
			timestamps, status = d.conn.SelectFAScriptTimestamps(book, chap)
			if status.IsErr {
				return status
			}
			if len(timestamps) > 0 {
				var bibleFileId int64
				bibleFileId, status = d.dbpConn.SelectFileId(filesetId,
					timestamps[0].BookId, timestamps[0].ChapterNum)
				if status.IsErr {
					return status
				}
				status = d.dbpConn.UpdateTimestamps(bibleFileId, timestamps)
				if status.IsErr {
					return status
				}
			}
		}
	}
	return status
}
