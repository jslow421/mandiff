package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func getRawText(ctx context.Context, client *textract.Client, jobID string) (string, error) {
	var rawText string
	nextToken := "" // Start with no token

	for {
		input := &textract.GetDocumentAnalysisInput{
			JobId:     aws.String(jobID),
			NextToken: aws.String(nextToken), // Pass NextToken if available
		}

		result, err := client.GetDocumentAnalysis(ctx, input)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Print detailed error information
				log.Println("AWS error code:", awsErr.Code(), "message:", awsErr.Message())
			} else {
				// Handle other types of errors
				log.Println(err.Error())
			}
			return "", err
		}

		for _, block := range result.Blocks {
			if block.BlockType == "LINE" {
				rawText += *block.Text + "\n" // Add newline if needed
			}
		}

		// Check for pagination
		nextToken = aws.ToString(result.NextToken)
		if nextToken == "" {
			break // No more pages
		}
	}

	return rawText, nil
}

func handleTextractCompletion(jobId string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	textractClient := textract.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	result, err := textractClient.GetDocumentAnalysis(context.TODO(), &textract.GetDocumentAnalysisInput{
		JobId: aws.String(jobId),
	})
	if err != nil {
		return fmt.Errorf("failed to get Textract results: %v", err)
	}

	var textBuilder strings.Builder
	for _, block := range result.Blocks {
		if block.BlockType == types.BlockTypeLine {
			textBuilder.WriteString(*block.Text)
			textBuilder.WriteString("\n")
		}
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("COMPLETE_BUCKET")),
		Key:         aws.String(jobId + ".txt"), // Assuming JobId is unique
		Body:        strings.NewReader(textBuilder.String()),
		ContentType: aws.String("text/plain"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload text file to S3: %v", err)
	}

	return nil
}

func handler() {
	fmt.Println("Processing documents...")

	// Replace placeholders
	jobID := "c7350199d357f286facbe43b29fdfd9acce301ee25f51d865836fd4d4deee8b2"
	bucketName := os.Getenv("COMPLETE_BUCKET")
	outputKey := "output.txt"

	// Load AWS credentials and create clients
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		os.Exit(1)
	}
	textractClient := textract.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	// Get the raw text
	rawText, err := getRawText(context.TODO(), textractClient, jobID)
	if err != nil {
		fmt.Println("Error getting raw text:", err)
		os.Exit(1)
	}

	// Upload to S3
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(outputKey),
		Body:   io.NopCloser(strings.NewReader(rawText)), // io.NopCloser for minimal memory usage
	})
	if err != nil {
		fmt.Println("Error uploading to S3:", err)
		os.Exit(1)
	}

	fmt.Println("Raw text stored in S3:", bucketName, "/", outputKey)

}

func main() {
	log.Println("Starting handler...")
	runtime.Start(handler)
}
