package input

import (
	"archive/zip"
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"dataset/request"
	"io"
	"os"
	"path/filepath"
	"sort"
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

func Unzip(ctx context.Context, files []InputFile) ([]InputFile, dataset.Status) {
	var results []InputFile
	var status dataset.Status
	if len(files) == 1 && filepath.Ext(files[0].Filename) == `.zip` {
		r, err := zip.OpenReader(files[0].FilePath())
		if err != nil {
			status = log.Error(ctx, 500, err, `Error unzipping file`)
			return results, status
		}
		defer r.Close()
		dest := files[0].Directory
		for _, f := range r.File {
			if f.FileInfo().Name()[0] == '.' {
				continue
			}
			rc, err2 := f.Open()
			if err2 != nil {
				status = log.Error(ctx, 500, err2, `Error reading zip file`)
				return results, status
			}
			defer rc.Close()
			if f.FileInfo().IsDir() {
				dest = filepath.Join(dest, f.FileInfo().Name())
				os.MkdirAll(dest, f.Mode())
				continue
			}
			path := filepath.Join(dest, f.FileInfo().Name())
			outFile, err3 := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err3 != nil {
				status = log.Error(ctx, 500, err3, `Error opening file to unzip into`)
				return results, status
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, rc)
			_ = rc.Close()
			if err != nil {
				status = log.Error(ctx, 500, err, `Error copying file during unzip`)
				return results, status
			}
			results = append(results, InputFile{Filename: f.FileInfo().Name(), Directory: dest})
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Filename < results[j].Filename
		})
	} else {
		results = files
	}
	return results, status
}

// SetMediaType function looks at names and sets the Media Type
func SetMediaType(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	fN := file.Filename
	if strings.HasSuffix(fN, `_ET`) || strings.HasSuffix(fN, `_ET.json`) {
		file.MediaType = request.TextPlainEdit
	} else if strings.HasSuffix(fN, `usx`) {
		file.MediaType = request.TextUSXEdit
	} else if (fN[0] == 'A' || fN[0] == 'B') && (fN[1] >= '0' && fN[1] <= '9') {
		file.MediaType = request.Audio
	} else if strings.HasSuffix(fN, `ST.xlsx`) || strings.HasSuffix(fN, `.xlsm`) {
		file.MediaType = request.TextScript
	} else if strings.HasSuffix(fN, `.csv`) {
		file.MediaType = request.TextCSV
	} else if (fN[0] == 'N' || fN[0] == 'O' || fN[0] == 'P') && fN[1] == '1' && fN[2] == '_' {
		file.MediaType = request.Audio
	} else if (fN[0] == 'N' || fN[0] == 'O' || fN[0] == 'P') && fN[1] == '2' && fN[2] == '_' {
		file.MediaType = request.AudioDrama
	} else {
		parts := strings.Split(fN, `_`)
		if len(parts) > 2 {
			ord := parts[1]
			if (ord[0] == 'A' || ord[0] == 'B') && (ord[1] >= '0' && ord[1] <= '9') {
				file.MediaType = request.Audio
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
	if file.MediaType == request.TextPlain || file.MediaType == request.TextPlainEdit ||
		file.MediaType == request.TextCSV {
		file.MediaId = strings.Split(file.Filename, `.`)[0]
		test := file.Filename[6]
		if test == 'O' {
			file.Testament = `OT`
		} else if test == 'N' {
			file.Testament = `NT`
		}
		file.FileExt = filepath.Ext(file.Filename)
	} else if file.MediaType == request.TextUSXEdit {
		parts := strings.Split(file.Directory, `/`)
		file.MediaId = parts[len(parts)-1]
		var tmpBookId, tmpBookSeq string
		if len(file.Filename) == 10 {
			tmpBookId = file.Filename[3:6]
			tmpBookSeq = file.Filename[0:3]
		} else if len(file.Filename) == 7 {
			tmpBookId = file.Filename[0:3]
			tmpBookSeq = strconv.Itoa(db.BookSeqMap[tmpBookId])
		} else {
			return log.ErrorNoErr(ctx, 400, `USX files are expected in the format 001GEN.usx or GEN.usx`)
		}
		file.BookId, status = validateBookId(ctx, tmpBookId)
		if status.IsErr {
			return status
		}
		file.BookSeq = tmpBookSeq
		file.Testament = db.Testament(file.BookId)
		file.FileExt = filepath.Ext(file.Filename)
	} else if file.MediaType == request.TextScript {
		file.MediaId = strings.Split(file.Filename, `.`)[0]
		test := file.Filename[6]
		if test == 'O' {
			file.Testament = `OT`
		} else if test == 'N' {
			file.Testament = `NT`
		}
		file.FileExt = filepath.Ext(file.Filename)
	} else if file.MediaType == request.Audio || file.MediaType == request.AudioDrama {
		fN := file.Filename
		if strings.HasSuffix(fN, `VOX.mp3`) {
			status = ParseVOXAudioFilename(ctx, file)
		} else if (fN[0] == 'A' || fN[0] == 'B') && (fN[1] >= '0' && fN[1] <= '9') {
			status = ParseV2AudioFilename(ctx, file)
		} else {
			status = ParseV4AudioFilename(ctx, file)
		}
		if status.IsErr {
			return status
		}
	} else {
		status = log.ErrorNoErr(ctx, 400, `Type must be one of "text_plain", "text_plain_edit", "text_usx", "audio"`)
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
		file.BookId, status = validateBookId(ctx, parts[2])
		if status.IsErr {
			return status
		}
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

func ParseVOXAudioFilename(ctx context.Context, file *InputFile) dataset.Status {
	var status dataset.Status
	var err error
	parts := strings.Split(file.Filename, `_`)
	if len(parts) != 7 {
		return log.ErrorNoErr(ctx, 500, "A VOX. filename is expected to have 7 parts", parts)
	}
	if parts[0][0] == 'N' {
		file.Testament = `NT`
	} else if parts[0][0] == 'O' {
		file.Testament = `OT`
	} else if parts[0][0] == 'P' {
		file.Testament = `PT` // what should this really be
	} else {
		return log.ErrorNoErr(ctx, 500, "Unknown media type", parts[0])
	}
	drama := parts[0]
	langCode := parts[1]
	versionCode := parts[2]
	file.BookSeq = parts[3]
	file.BookId, status = validateBookId(ctx, parts[4])
	if status.IsErr {
		return status
	}
	file.Chapter, err = strconv.Atoi(parts[5])
	if err != nil {
		return log.Error(ctx, 500, err, `Error convert chapter to int`, parts[5])
	}
	file.MediaId = langCode + versionCode + drama + "DA"
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

func UpdateIdent(conn db.DBAdapter, ident *db.Ident, textFiles []InputFile, audioFiles []InputFile) dataset.Status {
	var status dataset.Status
	_ = updateIdentText(ident, textFiles)
	_ = updateIdentAudio(ident, audioFiles)
	status = conn.InsertReplaceIdent(*ident)
	return status
}

func updateIdentText(ident *db.Ident, files []InputFile) bool {
	var result = false
	if len(files) > 0 {
		f := files[0]
		if f.MediaType != request.TextNone {
			ident.TextSource = f.MediaType
		}
		result = true
		if f.Testament == `OT` {
			ident.TextOTId = f.MediaId
			ident.LanguageISO = strings.ToLower(ident.TextOTId[:3])
		} else if f.Testament == `NT` {
			ident.TextNTId = f.MediaId
			ident.LanguageISO = strings.ToLower(ident.TextNTId[:3])
		}
	}
	return result
}

func updateIdentAudio(ident *db.Ident, files []InputFile) bool {
	var result = false
	for _, f := range files {
		result = true
		if f.Testament == `OT` {
			ident.AudioOTId = f.MediaId
			ident.LanguageISO = strings.ToLower(ident.AudioOTId[:3])
		} else if f.Testament == `NT` {
			ident.AudioNTId = f.MediaId
			ident.LanguageISO = strings.ToLower(ident.AudioNTId[:3])
		}
	}
	return result
}

var corrections = map[string]string{
	"EZE": "EZK", // Ezekiel
	"JMS": "JAS", // James
	"JOE": "JOL", // Joel
	"NAH": "NAM", // Nahum
	"PRV": "PRO", // Proverbs
	"PSM": "PSA", // Psalms
	"SOS": "SNG", // Song of Solomon
	"TTL": "TIT", // Titus
	"TTS": "TIT"} // Titus

func validateBookId(ctx context.Context, bookId string) (string, dataset.Status) {
	var status dataset.Status
	corrected, found := corrections[bookId]
	if found {
		bookId = corrected
	}
	_, ok := db.BookChapterMap[bookId]
	if !ok {
		status = log.ErrorNoErr(ctx, 500, "BookId", bookId, "is not known. Corrections:", corrections)
	}
	return bookId, status
}

// Parse DBP4 Audio names
//{mediaid}_{A/B}{ordering}_{USFM book code}_{chapter start}[_{verse start}-{chapter stop}_{verse stop}].mp3|webm
//eg: ENGESVN2DA_B001_MAT_001.mp3  (full chapter)
//eg: IRUNLCP1DA_B013_1TH_001_001-001_010.mp3  (partial chapter, verses 1-10)
