package sqs_client

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSProcessor struct {
	QueueURL string
}

func initSQSClient() (*sqs.Client, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	// Create and return the SQS client
	return sqs.NewFromConfig(cfg), nil
}

func pollMessagesSQS(queueURL string, client *sqs.Client) {
	log.Printf("Starting SQS polling... to Queue : %v \n", queueURL)
	for {
		log.Println("Polling for messages...")
		output, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     10,
		})
		if err != nil {
			log.Printf("Failed to Receive  Message :  %v \n", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if len(output.Messages) == 0 {
			log.Printf("No  Messages  received  ...")
			time.Sleep(2 * time.Second)
			continue
		}
		log.Printf("Received  %d  messages  :", len(output.Messages))

		// Process   messages and get successfull  ones
		processor := GetEventProcessor()
		successReceipts := processor.ProcessMessages(output.Messages)
		for _, receipt := range successReceipts {
			_, err = client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: aws.String(receipt),
			})
			if err != nil {
				log.Printf("Error deleting message: %v\n", err)
			} else {
				log.Println("Message deleted successfully.")
			}
		}
		time.Sleep(2 * time.Second)
	}
}
func (p *SQSProcessor) StartPolling() {
	client, err := initSQSClient()
	pollMessagesSQS(p.QueueURL, client)
	if err != nil {
		log.Fatalf("Failed to initialize SQS client: %v", err)
	}

}
