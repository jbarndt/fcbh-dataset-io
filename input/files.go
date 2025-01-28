package input

import (
	"context"
	log "dataset/logger"
	"dataset/request"
)

func FileInput(ctx context.Context, path string, testament request.Testament) ([]InputFile, *log.Status) {
	var files []InputFile
	var status *log.Status
	files, status = Glob(ctx, path)
	if status != nil {
		return files, status
	}
	files, status = Unzip(ctx, files)
	if status != nil {
		return files, status
	}
	for i, _ := range files {
		status = SetMediaType(ctx, &files[i])
		if status != nil {
			return files, status
		}
		status = ParseFilenames(ctx, &files[i])
		if status != nil {
			return files, status
		}
	}
	inputFiles := PruneBooksByRequest(files, testament)
	return inputFiles, status
}
