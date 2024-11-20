package input

import (
	"context"
	"dataset"
	"dataset/request"
)

func FileInput(ctx context.Context, path string, testament request.Testament) ([]InputFile, dataset.Status) {
	var files []InputFile
	var status dataset.Status
	files, status = Glob(ctx, path)
	if status.IsErr {
		return files, status
	}
	files, status = Unzip(ctx, files)
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
