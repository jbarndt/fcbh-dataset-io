package input

import (
	"context"
	"dataset"
	log "dataset/logger"
	"dataset/request"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	bucket, prefix, glob, status := parseGlob(ctx, path)
	if status.IsErr {
		return files, status
	}
	list, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		status = log.Error(ctx, 400, err, `Failed to list AWSS3 objects`)
		return files, status
	}
	bibleId, mediaId := findBibleIdMediaId(prefix)
	directory := filepath.Join(os.Getenv(`FCBH_DATASET_FILES`), bibleId, mediaId)
	status = EnsureDirectory(ctx, directory)
	for _, object := range list.Contents {
		objKey := aws.ToString(object.Key)
		if glob == nil || glob.MatchString(objKey) {
			var inFile InputFile
			inFile.Directory = directory
			inFile.Filename = filepath.Base(objKey)
			files = append(files, inFile)
			filePath := inFile.FilePath()
			fileInfo, stErr := os.Stat(filePath)
			if os.IsNotExist(stErr) || fileInfo.Size() != *object.Size {
				log.Info(ctx, `Downloading file`, objKey)
				response, getErr := client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(objKey),
				})
				if getErr != nil {
					log.Warn(ctx, getErr, `Failed to get object`, object.Key)
				}
				file, filErr := os.Create(filePath)
				if filErr != nil {
					status = log.Error(ctx, 400, filErr, `Failed to create file`, filePath)
					return files, status
				}
				_, copErr := io.Copy(file, response.Body)
				if err != nil {
					status = log.Error(ctx, 400, copErr, `Failed to copy object`, object.Key)
					return files, status
				}
				err = response.Body.Close()
				err = file.Close()
			}
		}
	}
	files, status = Unzip(ctx, files)
	if status.IsErr {
		return files, status
	}
	for i := range files {
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

func parseGlob(ctx context.Context, globKey string) (string, string, *regexp.Regexp, dataset.Status) {
	var bucket string
	var prefix string
	var regex *regexp.Regexp
	var status dataset.Status
	if strings.HasPrefix(globKey, `s3://`) {
		globKey = globKey[5:]
	} else if strings.HasPrefix(globKey, `s3:/`) {
		globKey = globKey[4:]
	}
	firstSlash := strings.Index(globKey, `/`)
	if firstSlash >= 0 {
		bucket = globKey[:firstSlash]
		prefix = globKey[firstSlash+1:]
		regex = nil
	}
	lastSlash := strings.LastIndex(globKey, `/`)
	if lastSlash >= 0 {
		glob := globKey[lastSlash+1:]
		if strings.Contains(glob, `*`) {
			prefix = globKey[firstSlash+1 : lastSlash+1]
			regex, status = globPattern(ctx, glob)
			if status.IsErr {
				return bucket, prefix, regex, status
			}
		}
	}
	return bucket, prefix, regex, status
}

func globPattern(ctx context.Context, glob string) (*regexp.Regexp, dataset.Status) {
	var regex *regexp.Regexp
	var status dataset.Status
	var err error
	glob = strings.Replace(glob, `.`, `\.`, -1)
	glob = strings.Replace(glob, `*`, `.`, -1)
	glob += `$`
	regex, err = regexp.Compile(glob)
	if err != nil {
		status = log.Error(ctx, 400, err, `Failed to compile glob pattern on AWS input`)
	}
	return regex, status
}

func findBibleIdMediaId(prefix string) (string, string) {
	var bibleId string
	var mediaId string
	parts := strings.Split(prefix, `/`)
	pos := len(parts) - 1
	for {
		if parts[pos] != `` {
			mediaId = parts[pos]
			bibleId = parts[pos-1]
			break
		}
		pos--
	}
	return bibleId, mediaId
}
