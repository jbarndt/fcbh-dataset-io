package controller

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"dataset"
	log "dataset/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	inputFolder  = `input/`
	sucessFolder = `success/`
	failedFolder = `failed/`
)

// main is the entry point for the ECS Task invocation.
// it expects to receive one S3 object (which is a yaml file); it will invoke the main service entry point, then exit
func main() {
	ctx := context.WithValue(context.Background(), `runType`, `queue`)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err, "aws LoadDefaultConfig Failed In Queue Main")
		os.Exit(1)
	}
	client := s3.NewFromConfig(cfg)

	// An Eventbridge event will set two environment variables containing the bucket and key to be retrieved
	// Permission to retrieve the object has been granted to the ECS task via terraform
	// values: S3_BUCKET and S3_KEY
	bucket := os.Getenv("S3_BUCKET")
	key := os.Getenv("S3_KEY")

	getInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	
	status dataset.Status
	object, err := client.GetObject(ctx, getInput)
	if err != nil {
		status = log.Error(ctx, 500, err, "Error Getting Object in Queue Input Folder")
		return
	}
	var content []byte

	content, err = io.ReadAll(object.Body)
	_ = object.Body.Close()
	if err != nil {
		status = log.Error(ctx, 500, err, "Error reading yaml file from Queue Input Folder.")
	}
	control := controller.NewController(ctx, object)
	_, status = control.ProcessV2() // calls main entry point
	var folder string
	if status.IsErr {
		folder = failedFolder
	} else {
		folder = sucessFolder
	}

	// first := true
	// for {
	// 	bucketName := os.Getenv("FCBH_DATASET_QUEUE")
	// 	object, key, status := getOldestObject(ctx, client, bucketName)
	// 	if first && status.IsErr {
	// 		_, _ = fmt.Fprintln(os.Stderr, err, "Reading First Input Failed In Queue Main")
	// 		os.Exit(1)
	// 	}
	// 	first = false
	// 	if !status.IsErr && object != nil {
	// 		control := controller.NewController(ctx, object)
	// 		_, status = control.ProcessV2() // calls main entry point
	// 		var folder string
	// 		if status.IsErr {
	// 			folder = failedFolder
	// 		} else {
	// 			folder = sucessFolder
	// 		}
	// 		_ = moveOnCompletion(ctx, client, bucketName, key, folder)
	// 	}
	// 	time.Sleep(time.Second * 10)
	// }
}

func getOldestObject(ctx context.Context, client *s3.Client, bucket string) ([]byte, string, dataset.Status) {
	var content []byte
	var key string
	var status dataset.Status
	inFolder := inputFolder
	if runtime.GOOS == "darwin" {
		inFolder = "input_test/"
	}
	input := &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: aws.String(inFolder),
	}
	result, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		status = log.Error(ctx, 500, err, "Error Listing Objects in Queue Input Folder")
		return content, key, status
	}
	if len(result.Contents) == 0 {
		return content, key, status
	}
	sort.Slice(result.Contents, func(i, j int) bool {
		return result.Contents[i].LastModified.Before(*result.Contents[j].LastModified)
	})
	key = *result.Contents[0].Key
	getInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	object, err := client.GetObject(ctx, getInput)
	if err != nil {
		status = log.Error(ctx, 500, err, "Error Getting Object in Queue Input Folder")
		return content, key, status
	}
	content, err = io.ReadAll(object.Body)
	_ = object.Body.Close()
	if err != nil {
		status = log.Error(ctx, 500, err, "Error reading yaml file from Queue Input Folder.")
	}
	return content, key, status
}

func moveOnCompletion(ctx context.Context, client *s3.Client, bucket, key string, folder string) dataset.Status {
	var status dataset.Status
	source := bucket + "/" + key
	dateTime := time.Now().Local().Format("2006-01-02T15:04:05")
	target := folder + dateTime + "-" + strings.Split(key, "/")[1]
	_, err := client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &bucket,
		CopySource: &source,
		Key:        &target,
	})
	if err != nil {
		status = log.Error(ctx, 500, err, "Error Moving File to", folder, "Folder")
	}
	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		status = log.Error(ctx, 500, err, "Error Deleting Object in Queue Input Folder")
	}
	return status
}
