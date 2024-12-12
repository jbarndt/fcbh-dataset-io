package main

import (
	"context"
	"dataset/controller"
	log "dataset/logger"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"os"
	"time"
)

// This is experimental code that did not fit requirements.
// SQS was a bad fit because transactions that are accidentally or intentional read
// are put into non-accessible state for as long as the longest transaction might run.
// The code is left here in case a solution to this problem is found

func main2() {
	var ctx = context.WithValue(context.Background(), `runType`, `sqs`)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err, "aws LoadDefaultConfig Failed")
		os.Exit(1)
	}
	sqsClient := sqs.NewFromConfig(cfg)
	queueURL := os.Getenv("FCBH_SQS_QUEUE")
	var first = true
	for {
		response, err2 := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})
		if err2 != nil {
			log.Error(ctx, 500, err2, "Error Receiving Message from SQS Queue")
			if first {
				os.Exit(1)
			}
			continue
		}
		first = false
		for _, message := range response.Messages {
			changedMsg := &sqs.ChangeMessageVisibilityInput{
				QueueUrl:          &queueURL,
				ReceiptHandle:     message.ReceiptHandle,
				VisibilityTimeout: *aws.Int32(43200),
			}
			_, err = sqsClient.ChangeMessageVisibility(ctx, changedMsg)
			if err != nil {
				log.Error(ctx, 500, err, "Error Changing Message Visibility")
			}
			yamlRequest := *message.Body
			var control = controller.NewController(ctx, []byte(yamlRequest))
			_, status := control.Process()
			if status.IsErr {
				err3 := sendToFailedQueue(ctx, sqsClient, yamlRequest, status.String())
				log.Error(ctx, 500, err3, "Error Writing Message to Failed SQS Queue")
			}
			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Error(ctx, 500, err, "Error Deleting SQS Message")
			}
		}
		time.Sleep(time.Second)
	}
}

func sendToFailedQueue(ctx context.Context, client *sqs.Client, jobData string, errorMsg string) error {
	queueURL := os.Getenv("FCBH_FAILED_QUEUE")
	_, err := client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &jobData,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"error": {
				DataType:    aws.String("String"),
				StringValue: aws.String(errorMsg),
			},
			"timestamp": {
				DataType:    aws.String("String"),
				StringValue: aws.String(time.Now().String()),
			},
		},
	})
	return err
}
