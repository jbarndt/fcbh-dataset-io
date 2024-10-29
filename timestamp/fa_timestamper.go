package timestamp

import (
	"context"
	"dataset"
	"dataset/db"
)

type FATimeStamper struct {
	ctx  context.Context
	conn db.DBAdapter
}

func NewFATimeStamper(ctx context.Context, conn db.DBAdapter) FATimeStamper {
	var f FATimeStamper
	f.ctx = ctx
	f.conn = conn
	return f
}

func (f *FATimeStamper) GetTimestamps(bookId string, chapterNum int) ([]db.Audio, dataset.Status) {
	var result []db.Audio
	var status dataset.Status
	result, status = f.conn.SelectFAScriptTimestamps(bookId, chapterNum)
	return result, status
}
