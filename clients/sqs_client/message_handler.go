package sqs_client

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/BRIZINGR007/go-service-utils/structs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type EventProcessor struct {
	EventHandlers map[string]func(json.RawMessage) error
}

var (
	instance *EventProcessor
	once     sync.Once
)

func GetEventProcessor() *EventProcessor {
	once.Do(func() {
		instance = &EventProcessor{}
	})
	return instance
}

func (p *EventProcessor) HandleEvents(messageBody structs.MessageBody) error {
	handler, exists := p.EventHandlers[messageBody.Event]
	if !exists {
		log.Printf("no handler registered for event: %s", messageBody.Event)
		return nil
	}

	return handler(messageBody.Payload)
}

func (p *EventProcessor) ProcessMessages(messages []types.Message) []string {
	var successReceipts []string
	for _, msg := range messages {
		log.Printf("Processing message ID: %s, Body: %s", *msg.MessageId, *msg.Body)
		var parsedBody structs.MessageBody
		if err := json.Unmarshal([]byte(*msg.Body), &parsedBody); err != nil {
			log.Printf("Failed to parse message body: %v", err)
			continue
		}
		if err := p.HandleEvents(parsedBody); err != nil {
			log.Printf("Error handling event: %v", err)
			continue
		}
		if msg.ReceiptHandle != nil {
			successReceipts = append(successReceipts, *msg.ReceiptHandle)
		}
	}
	return successReceipts
}
