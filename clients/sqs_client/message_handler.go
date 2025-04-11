package sqs_client

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"

	"github.com/BRIZINGR007/go-service-utils/structs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type EventProcessor struct {
	EventHandlers map[string]interface{}
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
		log.Printf("no handler registered for event : %s", messageBody.Event)
		return nil
	}
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	if handlerType.Kind() != reflect.Func || handlerType.NumIn() != 1 {
		log.Printf("invalid handler signature for event: %s", messageBody.Event)
		return nil
	}

	// Get the type of the function's first parameter
	argType := handlerType.In(0)

	// Convert map[string]string payload to JSON
	payloadBytes, err := json.Marshal(messageBody.Payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return nil
	}

	// Create a new instance of the argument type and unmarshal into it .
	// If we fail to verify the  payload then we delete the message from the  queue  .
	argValue := reflect.New(argType).Interface()
	if err := json.Unmarshal(payloadBytes, argValue); err != nil {
		log.Printf("failed to unmarshal payload into handler argument: %v", err)
		return nil
	}

	// Call the handler with the argument (dereference pointer if needed)
	reflectArg := reflect.ValueOf(argValue)
	if reflectArg.Kind() == reflect.Ptr {
		reflectArg = reflectArg.Elem()
	}

	returnValues := handlerValue.Call([]reflect.Value{reflectArg})

	if len(returnValues) > 0 && !returnValues[0].IsNil() {
		if returnValues[0].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return returnValues[0].Interface().(error)
		}
	}
	return nil

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
