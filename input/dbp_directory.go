package input

import (
	"context"
	"dataset/decode_yaml/request"
	log "dataset/logger"
	"os"
	"path/filepath"
)

// DBPDirectory 1. Assign pattern for OT, NT.  2. Glob files.  3. Assign book/chapter & Prune
func DBPDirectory(ctx context.Context, bibleId string, fsType request.MediaType, otFileset string,
	ntFileset string, testament request.Testament) ([]InputFile, *log.Status) {
	var results []InputFile
	var files []InputFile
	var status *log.Status
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
		if status != nil {
			return results, status
		}
		results = append(results, files...)
	}
	return results, status
}

func Directory(ctx context.Context, bibleId string, fsType request.MediaType, filesetId string, tType string,
	testament request.Testament) ([]InputFile, *log.Status) {
	var status *log.Status
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
	if status != nil {
		return files, status
	}
	for i, _ := range files {
		files[i].MediaType = fsType
		status = ParseFilenames(ctx, &files[i])
		if status != nil {
			return files, status
		}
	}
	inputFiles := PruneBooksByRequest(files, testament)
	return inputFiles, status
}
