package sqs_client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/BRIZINGR007/go-service-utils/structs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func PostMessageNonFIFO(queueURL string, payload *structs.MessageBody) error {
	client, err := InitSQSClient()
	if err != nil {
		return fmt.Errorf("failed to initialize SQS client: %w", err)
	}

	messageBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message body: %w", err)
	}

	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to queue %s: %w", queueURL, err)
	}
	log.Println("Message sent successfully to queue:", queueURL)
	return nil
}
