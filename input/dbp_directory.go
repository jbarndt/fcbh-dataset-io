package input

import (
	"context"
	"dataset"
	"dataset/request"
	"os"
	"path/filepath"
)

// DBPDirectory 1. Assign pattern for OT, NT.  2. Glob files.  3. Assign book/chapter & Prune
func DBPDirectory(ctx context.Context, bibleId string, fsType request.MediaType, otFileset string,
	ntFileset string, testament request.Testament) ([]InputFile, dataset.Status) {
	var results []InputFile
	var files []InputFile
	var status dataset.Status
	type run struct {
		filesetId string
		tType     string
	}
	var runs []run
	if otFileset != `` {
		runs = append(runs, run{filesetId: otFileset, tType: `OT`})
	}
	if ntFileset != `` {
		runs = append(runs, run{filesetId: ntFileset, tType: `NT`})
	}
	for _, r := range runs {
		files, status = Directory(ctx, bibleId, fsType, r.filesetId, r.tType, testament)
		if status.IsErr {
			return results, status
		}
		results = append(results, files...)
	}
	return results, status
}

func Directory(ctx context.Context, bibleId string, fsType request.MediaType, filesetId string, tType string,
	testament request.Testament) ([]InputFile, dataset.Status) {
	var status dataset.Status
	var directory string
	var search string
	if fsType == request.TextPlain || fsType == request.TextPlainEdit {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
		search = filepath.Join(directory, filesetId+".json")
	} else if fsType == request.TextUSXEdit {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, filesetId)
		search = filepath.Join(directory, "*.usx")
	} else if fsType == request.Audio || fsType == request.AudioDrama {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, filesetId)
		if tType == `OT` {
			search = filepath.Join(directory, "A*.*")
		} else {
			search = filepath.Join(directory, "B*.*")
		}
	}
	//fmt.Println("search:", tType, search)
	var files []InputFile
	files, status = Glob(ctx, search)
	if status.IsErr {
		return files, status
	}
	for i, _ := range files {
		files[i].MediaType = fsType
		status = ParseFilenames(ctx, &files[i])
		if status.IsErr {
			return files, status
		}
	}
	inputFiles := PruneBooksByRequest(files, testament)
	return inputFiles, status
}
