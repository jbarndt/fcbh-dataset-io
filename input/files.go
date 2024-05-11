package input

import (
	"context"
	"dataset"
	log "dataset/logger"
	"dataset/request"
	"os"
	"path/filepath"
)

func FileInput(ctx context.Context, path string, testament request.Testament) ([]InputFile, dataset.Status) {
	var files []InputFile
	var status dataset.Status
	files, status = Glob(ctx, path)
	if status.IsErr {
		return files, status
	}
	for i, _ := range files {
		status = SetMediaType(ctx, &files[i])
		if status.IsErr {
			return files, status
		}
		status = ParseFilenames(ctx, &files[i])
		if status.IsErr {
			return files, status
		}
	}
	inputFiles := PruneBooksByRequest(files, testament)
	return inputFiles, status
}

// PostInput corrects the name of the uploaded file to have the filename
// That was put in the request in the line POST
func PostInput(ctx context.Context, filePath string, post string, testament request.Testament) ([]InputFile, dataset.Status) {
	var files []InputFile
	var status dataset.Status
	directory := filepath.Dir(filePath)
	correctName := filepath.Join(directory, post)
	err := os.Rename(filePath, correctName)
	if err != nil {
		status = log.Error(ctx, 500, err, `Could not rename posted file`)
		return files, status
	}
	return FileInput(ctx, correctName, testament)
}
