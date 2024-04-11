package fetch

import (
	"context"
	"dataset/request"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPlainText(t *testing.T) {
	ctx := context.Background()
	var req request.Request
	req.AudioData.NoAudio = true
	req.TextData.BibleBrain.TextPlainEdit = true
	req.Testament.NT = true
	req.Testament.OT = true
	client := NewAPIDBPClient(ctx, `ENGWEB`)
	info, status := client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
		if ok {
			if len(info.AudioFilesets) > 0 {
				t.Error(`Should have found no audio filesets.`)
			}
			if len(info.TextFilesets) != 2 {
				t.Error(`Should have found two text filesets`)
			}
			err := deleteFile(info.BibleId, info.TextFilesets[1])
			if err != nil {
				t.Error(`Did not delete file.`)
			}
			download := NewAPIDownloadClient(ctx, info.BibleId)
			status := download.Download(info)
			if status.IsErr {
				t.Error("Unexpected Error", status.Message)
			}
			count, err := countFiles(info.BibleId, info.TextFilesets[1])
			if count != 2 {
				t.Error("Two text files are expected")
			}
			if err != nil {
				t.Error(`Download err was not expected`, err)
			}
		}
	}
}

func TestUSXDownload(t *testing.T) {
	ctx := context.Background()
	var req request.Request
	req.AudioData.NoAudio = true
	req.TextData.BibleBrain.TextUSXEdit = true
	req.Testament.NT = true
	client := NewAPIDBPClient(ctx, `ENGWEB`)
	info, status := client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
		if ok {
			if len(info.AudioFilesets) > 0 {
				t.Error(`Should have found no audio filesets.`)
			}
			if len(info.TextFilesets) != 1 {
				t.Error(`Should have found one text fileset`)
			}
			//fmt.Println(info.TextFilesets)
			err := deleteFile(info.BibleId, info.TextFilesets[0])
			if err != nil {
				t.Error(`Did not delete file.`)
			}
			download := NewAPIDownloadClient(ctx, info.BibleId)
			status := download.Download(info)
			if status.IsErr {
				t.Error(`Download Err is unexpected`, status.Message)
			}
			count, err := countFiles(info.BibleId, info.TextFilesets[0])
			if count != 27 {
				t.Error("27 books in NT are expected, found:", count)
			}
			if err != nil {
				t.Error(`Download err was not expected`, err)
			}
		}
	}
}

func TestAudioDownload(t *testing.T) {
	ctx := context.Background()
	var req request.Request
	req.AudioData.BibleBrain.MP3_64 = true
	req.TextData.NoText = true
	req.Testament.NT = true
	client := NewAPIDBPClient(ctx, `ENGWEB`)
	info, status := client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
		if ok {
			if len(info.AudioFilesets) != 1 {
				t.Error(`Should have found no audio filesets.`, len(info.AudioFilesets))
			}
			if len(info.TextFilesets) > 0 {
				t.Error(`Should have found no text fileset`)
			}
			//fmt.Println(info.TextFilesets)
			err := deleteFile(info.BibleId, info.AudioFilesets[0])
			if err != nil {
				t.Error(`Did not delete file.`)
			}
			download := NewAPIDownloadClient(ctx, info.BibleId)
			status := download.Download(info)
			if status.IsErr {
				t.Error(`Download Err is unexpected`, status.Message)
			}
			count, err := countFiles(info.BibleId, info.AudioFilesets[0])
			if count != 260 {
				t.Error("260 chapters in NT are expected, found:", count)
			}
			if err != nil {
				t.Error(`Download err was not expected`, err)
			}
		}
	}
}

func Test403Error(t *testing.T) {
	ctx := context.Background()
	var req request.Request
	req.AudioData.BibleBrain.MP3_64 = true
	req.TextData.NoText = true
	req.Testament.NT = true
	client := NewAPIDBPClient(ctx, `ENGESV`)
	info, status := client.BibleInfo()
	if !status.IsErr {
		ok := client.FindFilesets(&info, req.AudioData.BibleBrain, req.TextData.BibleBrain, req.Testament)
		if ok {
			if len(info.AudioFilesets) != 1 {
				t.Error(`Should have found no audio filesets.`, len(info.AudioFilesets))
			}
			if len(info.TextFilesets) > 0 {
				t.Error(`Should have found no text fileset`)
			}
			//fmt.Println(info.TextFilesets)
			err := deleteFile(info.BibleId, info.AudioFilesets[0])
			if err != nil {
				t.Error(`Did not delete file.`)
			}
			download := NewAPIDownloadClient(ctx, info.BibleId)
			status := download.Download(info)
			if status.Status != 403 {
				t.Error(`Download 403 error is unexpected`, status)
			}
		}
	}
}

func deleteFile(bibleId string, fs FilesetType) error {
	if fs.Type == `text_plain` {
		filePath := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, fs.Id+`.json`)
		fmt.Println(`Delete`, filePath)
		err := os.Remove(filePath)
		return err
	} else {
		filePath := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, fs.Id)
		files, err := os.ReadDir(filePath)
		if err != nil {
			os.IsNotExist(err)
		} else {
			panic(err)
		}
		for _, num := range []int{11, 7, 5, 3, 1} {
			if len(files) > num {
				delPath := filepath.Join(filePath, files[num].Name())
				fmt.Println("Delete", delPath)
				err := os.Remove(delPath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func countFiles(bibleId string, fs FilesetType) (int, error) {
	var filePath string
	if fs.Type == `text_plain` {
		filePath = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, `*.json`)
	} else if fs.Type == `text_usx` {
		filePath = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, fs.Id, `*.usx`)
	} else {
		filePath = filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, fs.Id, `*.`+fs.Container)
	}
	files, err := filepath.Glob(filePath)
	return len(files), err
}
