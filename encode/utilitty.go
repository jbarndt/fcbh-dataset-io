package encode

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadDirectory(ctx context.Context, bibleId string, filesetId string) ([]string, dataset.Status) {
	var results []string
	var status dataset.Status
	var directory = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, filesetId)
	var files, err = os.ReadDir(directory)
	if err != nil {
		status = log.Error(ctx, 500, err, `Error reading audio file directory.`)
		return results, status
	}
	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), `.`) {
			fullPath := filepath.Join(directory, file.Name())
			results = append(results, fullPath)
		}
	}
	return results, status
}

func ParseFilename(ctx context.Context, filePath string) (string, int, dataset.Status) {
	var bookId string
	var chapterNum int
	var status dataset.Status
	filename := filepath.Base(filePath)
	chapter, err := strconv.Atoi(filename[6:8])
	if err != nil {
		status = log.Error(ctx, 500, err, `Error convert chapter to int`, filename[6:8])
		return bookId, chapterNum, status
	}
	book := strings.Trim(filename[9:21], `_`)
	bookId = db.USFMBookId(ctx, book)
	return bookId, chapter, status
}
