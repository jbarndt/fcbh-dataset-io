package input

import (
	"context"
	"dataset"
)

func FileInput(ctx context.Context, path string) ([]InputFile, dataset.Status) {
	var files []InputFile
	var status dataset.Status
	files, status = Glob(ctx, path)
	if status.IsErr {
		return files, status
	}
	for _, file := range files {
		status = SetMediaType(ctx, &file)
		if status.IsErr {
			return files, status
		}
	}
	return files, status
}
