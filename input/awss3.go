package input

import (
	"context"
	"dataset"
	log "dataset/logger"
	"dataset/request"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// https://aws.github.io/aws-sdk-go-v2/docs/

// AWSS3Input is given a path prefix, that it uses to identify files.
// Saves each file found to disk, and returns an array of input files
func AWSS3Input(ctx context.Context, path string, testament request.Testament) ([]InputFile, dataset.Status) {
	var files []InputFile
	var status dataset.Status
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		status = log.Error(ctx, 400, err, `Failed to load AWS configuration`)
		return files, status
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "us-west-2"
	})
	//if !strings.HasSuffix(path, `/`) {
	//	path = path + `/`
	//}
	pathParts := strings.Split(path, `/`)
	bucket := pathParts[2]
	prefix := strings.Join(pathParts[3:], `/`)
	list, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		status = log.Error(ctx, 400, err, `Failed to list AWSS3 objects`)
		return files, status
	}
	mediaId := pathParts[len(pathParts)-2]
	bibleId := pathParts[len(pathParts)-3]
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, mediaId)
	status = EnsureDirectory(ctx, directory)
	for _, object := range list.Contents {
		log.Info(ctx, "Key=", aws.ToString(object.Key), "size=", *object.Size)
		var inFile InputFile
		inFile.Directory = directory
		inFile.Filename = filepath.Base(aws.ToString(object.Key))
		files = append(files, inFile)
		filePath := inFile.FilePath()
		fileInfo, err := os.Stat(filePath)
		if os.IsNotExist(err) || fileInfo.Size() != *object.Size {
			response, err := client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(*object.Key),
			})
			if err != nil {
				log.Warn(ctx, err, `Failed to get object`, object.Key)
			}
			file, err := os.Create(filePath)
			if err != nil {
				status = log.Error(ctx, 400, err, `Failed to create file`, filePath)
				return files, status
			}
			var count int64
			count, err = io.Copy(file, response.Body)
			if err != nil {
				status = log.Error(ctx, 400, err, `Failed to copy object`, object.Key)
				return files, status
			}
			fmt.Println("size downloaded", count)
		}
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

func EnsureDirectory(ctx context.Context, directory string) dataset.Status {
	var status dataset.Status
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		err2 := os.MkdirAll(directory, os.ModePerm)
		if err2 != nil {
			status = log.Error(ctx, 400, err2, `Failed to create directory to download files`)
		}
	} else if err != nil {
		status = log.Error(ctx, 400, err, `Failed to stat directory`)
	}
	return status
}
