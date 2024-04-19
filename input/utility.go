package input

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"dataset/request"
	"path/filepath"
	"strconv"
	"strings"
)

func Glob(ctx context.Context, search string) ([]InputFile, dataset.Status) {
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
			input.Directory = filepath.Dir(file)
			input.Filename = filepath.Base(file)
			results = append(results, input)
		}
	}
	return results, status
}

// SetMediaType function looks at names and sets the Media Type
func SetMediaType(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	fN := file.Filename
	if strings.HasSuffix(fN, `_ET`) || strings.HasSuffix(fN, `_ET.json`) {
		file.MediaType = `text_plain`
	} else if strings.HasSuffix(fN, `usx`) {
		file.MediaType = `text_usx`
	} else if (fN[0] == 'A' || fN[0] == 'B') && (fN[1] >= '0' && fN[1] <= '9') {
		file.MediaType = `audio`
	} else {
		parts := strings.Split(fN, `_`)
		if len(parts) > 2 {
			ord := parts[1]
			if (ord[0] == 'A' || ord[0] == 'B') && (ord[1] >= '0' && ord[1] <= '9') {
				file.MediaType = `audio`
			}
		}
	}
	if file.MediaType == `` {
		status = log.ErrorNoErr(ctx, 400, `Could not identify media type from filename`)
	}
	return status
}

func ParseFilenames(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	if file.MediaType == `text_plain` {
		file.MediaId = strings.Split(file.Filename, `.`)[0]
		test := file.Filename[6]
		if test == 'O' {
			file.Testament = `OT`
		} else if test == 'N' {
			file.Testament = `NT`
		}
		file.FileExt = filepath.Ext(file.Filename)
	} else if file.MediaType == `text_usx` {
		parts := strings.Split(file.Directory, `/`)
		file.MediaId = parts[len(parts)-1]
		file.BookId = file.Filename[3:6]
		file.BookSeq = file.Filename[0:3]
		file.Testament = db.Testament(file.BookId)
		file.FileExt = filepath.Ext(file.Filename)
	} else if file.MediaType == `audio` || file.MediaType == `audio_drama` {
		fN := file.Filename
		if (fN[0] == 'A' || fN[0] == 'B') && (fN[1] >= '0' && fN[1] <= '9') {
			status = ParseV2AudioFilename(ctx, file)
		} else {
			status = ParseV4AudioFilename(ctx, file)
		}
		if status.IsErr {
			return status
		}
	} else {
		status = log.ErrorNoErr(ctx, 400, `Type must be one of "text_plain", "text_usx", "audio"`)
	}
	return status
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

func PruneBooksByRequest(files []InputFile, testament request.Testament) []InputFile {
	var results []InputFile
	for _, f := range files {
		if testament.Has(f.Testament, f.BookId) || f.BookId == `` {
			results = append(results, f)
		}
	}
	return results
}

// Parse DBP4 Audio names
//{mediaid}_{A/B}{ordering}_{USFM book code}_{chapter start}[_{verse start}-{chapter stop}_{verse stop}].mp3|webm
//eg: ENGESVN2DA_B001_MAT_001.mp3  (full chapter)
//eg: IRUNLCP1DA_B013_1TH_001_001-001_010.mp3  (partial chapter, verses 1-10)
