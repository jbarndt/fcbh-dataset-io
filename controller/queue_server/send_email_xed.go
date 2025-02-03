package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	log "github.com/faithcomesbyhearing/fcbh-dataset-io/logger"
)

// This is just example code.  I need to have a verified email address to send from
// Or, I need to get a TXT entry entered into DNS.

func sendEmailSES(fromAddress, toAddress, subject, body string) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"), // Replace with your AWS region
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create SES client
	sesClient := ses.NewFromConfig(cfg)

	// Create email input
	input := &ses.SendEmailInput{
		Source: &fromAddress,
		Destination: &types.Destination{
			ToAddresses: []string{toAddress},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: &subject,
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: &body,
				},
			},
		},
	}

	// Send email
	_, err = sesClient.SendEmail(context.TODO(), input)
	return err
}

func main1() {
	fromEmail := "sender@yourdomain.com" // Must be verified in SES
	toEmail := "recipient@example.com"
	subject := "Test Email from AWS SES"
	body := "This is a test email sent using Amazon SES."

	err := sendEmailSES(fromEmail, toEmail, subject, body)
	if err != nil {
		log.Error(context.TODO(), 500, err, "Failed to send email")
	}
	fmt.Println("Email sent successfully!")
}
