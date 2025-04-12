package sqs_client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/BRIZINGR007/go-service-utils/structs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func PostMessageNonFIFO(queueURL string, payload *structs.MessageBody) {
	client, err := InitSQSClient()
	if err != nil {
		log.Fatalf("Failed to initialize SQS client: %v", err)
	}
	messageBody, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal message body: %v", err)
	}

	// Send the message
	_, err = client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: aws.String(string(messageBody)),
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Message sent successfully to queue:", queueURL)

}
