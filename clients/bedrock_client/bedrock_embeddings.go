package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func initBedrockRuntimeClient() (*bedrockruntime.Client, error) {
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
	return bedrockruntime.NewFromConfig(cfg), nil
}
func GenerateEmbeddings(inputText string, modelID string) ([]float32, error) {
	bedrockClient, err := initBedrockRuntimeClient()
	if err != nil {
		return nil, err
	}
	if modelID == "" {
		modelID = "amazon.titan-embed-text-v2:0"
	}
	payload := struct {
		InputText string `json:"inputText"`
	}{
		InputText: inputText,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Body:        payloadBytes,
	}

	// Invoke the model
	output, err := bedrockClient.InvokeModel(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke Bedrock model: %w", err)
	}

	// Parse the response
	var response struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.Unmarshal(output.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Embedding, nil
}
