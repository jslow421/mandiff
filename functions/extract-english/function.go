package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const outputFileKey = "output-english.txt"

func Handler(ctx context.Context, event events.S3Event) error {
	// Retrieve the .txt file key from the event
	//inputObjectKey := event.Records[0].S3.Object.Key
	// Temp hardcoded object key
	inputObjectKey := "output.txt"

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("error loading AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	comprehendClient := comprehend.NewFromConfig(cfg)

	// Fetch the .txt content from S3
	input := &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("SOURCE_BUCKET")),
		Key:    aws.String(inputObjectKey),
	}

	result, err := s3Client.GetObject(ctx, input)
	if err != nil {
		return fmt.Errorf("error reading from S3: %v", err)
	}
	defer result.Body.Close()

	textBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("error reading S3 object body: %v", err)
	}
	text := string(textBytes)

	englishText, err := detectEnglishBlocks(ctx, comprehendClient, text)
	if err != nil {
		return fmt.Errorf("error detecting English: %v", err)
	}

	// Write to S3
	output := &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("TARGET_BUCKET")),
		Key:    aws.String(outputFileKey),
		Body:   strings.NewReader(englishText),
	}

	_, err = s3Client.PutObject(ctx, output)
	if err != nil {
		return fmt.Errorf("error writing to S3: %v", err)
	}

	return nil
}

func detectEnglishBlocks(ctx context.Context, client *comprehend.Client, text string) (string, error) {
	textChunks := strings.Split(text, "\n")
	var englishBlocks []string

	for _, chunk := range textChunks {
		input := &comprehend.DetectDominantLanguageInput{
			Text: aws.String(chunk),
		}
		output, err := client.DetectDominantLanguage(ctx, input)
		if err != nil {
			return "", err
		}

		// Check if English is present at all
		englishDetected := false
		for _, lang := range output.Languages {
			if *lang.LanguageCode == "en" {
				englishDetected = true
				break
			}
		}

		if englishDetected {
			englishBlocks = append(englishBlocks, chunk)
		}
	}

	return strings.Join(englishBlocks, "\n"), nil
}

func main() {
	log.Println("Starting handler...")
	lambda.Start(Handler)
}
