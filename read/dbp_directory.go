package read

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"dataset/request"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type InputFile struct {
	MediaId    string
	Testament  string
	BookId     string // not used for text_plain
	BookSeq    string
	Chapter    int    // only used for audio
	Verse      string // used by OBT and Vessel
	ChapterEnd int
	VerseEnd   string
	Filename   string
	FileExt    string
	Directory  string
}

func (i *InputFile) FilePath() string {
	return filepath.Join(i.Directory, i.Filename)
}

// DBPDirectory 1. Assign pattern for OT, NT.  2. Glob files.  3. Assign book/chapter & Prune
func DBPDirectory(ctx context.Context, bibleId string, fsType string, otFileset string, ntFileset string,
	testament request.Testament) ([]InputFile, dataset.Status) {
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

func Directory(ctx context.Context, bibleId string, fsType string, filesetId string, tType string,
	testament request.Testament) ([]InputFile, dataset.Status) {
	var status dataset.Status
	var directory string
	var search string
	if fsType == `text_plain` {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId)
		search = filepath.Join(directory, filesetId+".json")
	} else if fsType == `text_usx` {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, filesetId)
		search = filepath.Join(directory, "*.usx")
	} else if fsType == `audio` {
		directory = filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, filesetId)
		if tType == `OT` {
			search = filepath.Join(directory, "A*.*")
		} else {
			search = filepath.Join(directory, "B*.*")
		}
	}
	//fmt.Println("search:", tType, search)
	var files []InputFile
	files, status = Glob(ctx, directory, search)
	if status.IsErr {
		return files, status
	}
	var inputFiles []InputFile
	if fsType == `text_plain` {
		inputFiles = files
	} else if fsType == `text_usx` {
		for _, file := range files {
			file.BookId = file.Filename[3:6]
			if testament.Has(tType, file.BookId) {
				file.FileExt = filepath.Ext(file.Filename)
				file.Testament = tType
				file.BookSeq = file.Filename[0:3]
				inputFiles = append(inputFiles, file)
			}
		}
	} else if fsType == `audio` {
		for _, file := range files {
			fN := file.Filename
			if (fN[0] == 'A' || fN[0] == 'B') && (fN[1] >= '0' && fN[1] <= '9') {
				status = ParseV2AudioFilename(ctx, &file)
			} else {
				status = ParseV4AudioFilename(ctx, &file)
			}
			if status.IsErr {
				return inputFiles, status
			}
			if testament.Has(tType, file.BookId) {
				inputFiles = append(inputFiles, file)
			}
		}
	} else {
		status = log.ErrorNoErr(ctx, 500, `Type must be one of "text_plain", "text_usx", "audio"`)
	}
	return inputFiles, status
}

func Glob(ctx context.Context, directory string, search string) ([]InputFile, dataset.Status) {
	var results []InputFile
	var status dataset.Status
	if search != `` {
		files, err := filepath.Glob(search)
		if err != nil {
			status = log.Error(ctx, 500, err, `Error reading directory`)
			return results, status
		}
		for _, file := range files {
			var input InputFile
			input.Directory = directory
			input.Filename = filepath.Base(file)
			results = append(results, input)
		}
	}
	return results, status
}

func ParseV2AudioFilename(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	var err error
	file.FileExt = filepath.Ext(file.Filename)
	filename := file.Filename[:len(file.Filename)-len(file.FileExt)]
	ab := filename[0]
	if ab == 'A' {
		file.Testament = `OT`
	} else if ab == 'B' {
		file.Testament = `NT`
	}
	seq := filename[1:4]
	file.BookSeq = strings.Trim(seq, `_`)
	file.Chapter, err = strconv.Atoi(file.Filename[6:8])
	if err != nil {
		return log.Error(ctx, 500, err, `Error convert chapter to int`, file.Filename[6:8])
	}
	book := strings.Trim(filename[9:21], `_`)
	file.BookId = db.USFMBookId(ctx, book)
	file.MediaId = filename[21:]
	return status
}

func ParseV4AudioFilename(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	var err error
	file.FileExt = filepath.Ext(file.Filename)
	filename := file.Filename[:len(file.Filename)-len(file.FileExt)]
	filename = strings.Replace(filename, `-`, `_`, -1)
	parts := strings.Split(filename, `_`)
	file.MediaId = parts[0]
	if len(parts) > 1 {
		ab := parts[1][0]
		if ab == 'A' {
			file.Testament = `OT`
		} else if ab == 'B' {
			file.Testament = `NT`
		}
		file.BookSeq = parts[1][1:]
	}
	if len(parts) > 2 {
		file.BookId = parts[2]
	}
	if len(parts) > 3 {
		file.Chapter, err = strconv.Atoi(parts[3])
		if err != nil {
			return log.Error(ctx, 500, err, `Error convert chapter to int`, parts[3])
		}
	}
	if len(parts) > 4 {
		file.Verse = parts[4]
	}
	if len(parts) > 5 {
		file.ChapterEnd, err = strconv.Atoi(parts[5])
		if err != nil {
			return log.Error(ctx, 500, err, `Error convert chapterEnd to int`, parts[5])
		}
	}
	if len(parts) > 6 {
		file.VerseEnd = parts[6]
	}
	return status
}

// Parse DBP4 Audio names
//{mediaid}_{A/B}{ordering}_{USFM book code}_{chapter start}[_{verse start}-{chapter stop}_{verse stop}].mp3|webm
//eg: ENGESVN2DA_B001_MAT_001.mp3  (full chapter)
//eg: IRUNLCP1DA_B013_1TH_001_001-001_010.mp3  (partial chapter, verses 1-10)
