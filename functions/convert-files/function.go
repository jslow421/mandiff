package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pemistahl/lingua-go"
)

func getRawText(ctx context.Context, client *textract.Client, jobID string) (string, error) {
	var rawText string
	nextToken := ""

	for {
		input := &textract.GetDocumentAnalysisInput{
			JobId: aws.String(jobID),
		}

		if nextToken != "" {
			input.NextToken = &nextToken
		}

		result, err := client.GetDocumentAnalysis(ctx, input)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				// Print detailed error information
				log.Println("AWS error code:", awsErr.Code(), "message:", awsErr.Message())
			} else {
				// Handle other types of errors
				log.Println("Some other error:", err.Error())
			}
			return "", err
		}

		for _, block := range result.Blocks {
			if block.BlockType == "LINE" {
				// Detect language
				value := *block.Text
				isEnglish, _ := checkIsEnglish(value)
				if isEnglish {
					rawText += value + "\n"
				}
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

func checkIsEnglish(line string) (bool, error) {
	var output string
	var confidence float64
	languages := []lingua.Language{
		lingua.English,
		lingua.French,
		lingua.German,
		lingua.Swedish,
		lingua.Italian,
		lingua.Spanish,
		lingua.Portuguese,
		lingua.Polish,
		lingua.Dutch,
		lingua.Danish,
		lingua.Finnish,
		lingua.Lithuanian,
	}

	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(languages...).Build()

	if language, exists := detector.DetectLanguageOf(line); exists {
		output = language.String()
		confidence = detector.ComputeLanguageConfidence(output, language)
	}

	isEnglish := (output == "English" && confidence >= 0.2)

	return isEnglish, nil
}

type DocumentEvent struct {
	JobId            string `json:"jobId"`
	InputBucketName  string `json:"inputBucketName"`
	OutputBucketName string `json:"outputBucketName"`
	OutputFileName   string `json:"outputFileName"`
}

func handler(ctx context.Context, event *DocumentEvent) (string, error) {
	// Replace placeholders
	jobID := event.JobId
	bucketName := os.Getenv("COMPLETE_BUCKET")
	outputKey := "output.txt"

	log.Println("Processing documents with job ID:", jobID)

	// Load AWS credentials and create clients
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading AWS config:", err)
		os.Exit(1)
	}
	textractClient := textract.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	// Get the raw text
	rawText, err := getRawText(context.TODO(), textractClient, jobID)
	if err != nil {
		log.Println("Error getting raw text:", err)
		os.Exit(1)
	}

	if rawText == "" {
		log.Println("No raw text found - this seems unlikely. Perhaps there are no English lines in the document?")
		os.Exit(1)
	}

	// Upload to S3
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(outputKey),
		Body:   io.NopCloser(strings.NewReader(rawText)),
	})
	if err != nil {
		fmt.Println("Error uploading to S3:", err)
		os.Exit(1)
	}

	log.Println("Raw text stored in S3:", bucketName, "/", outputKey)
	return "Success", nil
}

func main() {
	log.Println("Starting handler...")
	lambda.Start(handler)
}
