package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/pemistahl/lingua-go"
)

func getRawText(ctx context.Context, client *textract.Client, jobID string, shouldFilterLanguage bool) (string, error) {
	var rawText string
	nextToken := ""

	for {
		input := &textract.GetDocumentAnalysisInput{
			JobId: &jobID,
		}

		if nextToken != "" {
			input.NextToken = &nextToken
		}

		result, err := client.GetDocumentAnalysis(ctx, input)
		if err != nil {
			log.Fatal("Error getting document analysis:", err.Error())
			return "", err
		}

		for _, block := range result.Blocks {
			rawTextValue := getRawTextValue(block, shouldFilterLanguage)
			rawText += rawTextValue + "\n"
		}

		if nextToken == "" {
			break
		} else {
			nextToken = *result.NextToken
		}
	}

	return rawText, nil
}

func getRawTextValue(block types.Block, shouldFilterLanguage bool) string {
	var rawText string
	if block.BlockType == "LINE" {
		value := *block.Text

		if shouldFilterLanguage {
			isEnglish := checkIfTextIsLikelyEnglish(value)
			if isEnglish {
				rawText += value + "\n"
			}
		} else {
			rawText += value + "\n"
		}
	}

	return rawText
}

func checkIfTextIsLikelyEnglish(line string) bool {
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
		lingua.Turkish,
	}

	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(languages...).Build()

	if language, exists := detector.DetectLanguageOf(line); exists {
		output = language.String()
		confidence = detector.ComputeLanguageConfidence(output, language)
	}

	isEnglish := (output == "English" && confidence >= 0.2)

	return isEnglish
}

type DocumentEvent struct {
	JobId                string `json:"jobId"`
	InputBucketName      string `json:"inputBucketName"`
	OutputBucketName     string `json:"outputBucketName"`
	OutputFileName       string `json:"outputFileName"`
	ShouldFilterLanguage bool   `json:"shouldFilterLanguage"`
}

func handler(ctx context.Context, event *DocumentEvent) (string, error) {
	jobID := event.JobId
	bucketName := os.Getenv("COMPLETE_BUCKET")

	log.Println("Processing documents with job ID:", jobID)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println("Error loading AWS config:", err)
		os.Exit(1)
	}
	textractClient := textract.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	rawText, err := getRawText(context.TODO(), textractClient, jobID, event.ShouldFilterLanguage)
	if err != nil {
		log.Println("Error getting raw text:", err)
		os.Exit(1)
	}

	if rawText == "" {
		log.Println("No raw text found - this seems unlikely. Perhaps there are no English lines in the document?")
		os.Exit(1)
	}

	// Upload to S3
	fileName := event.OutputFileName + ".txt"
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   io.NopCloser(strings.NewReader(rawText)),
	})
	if err != nil {
		log.Println("Error uploading to S3:", err)
		os.Exit(1)
	}

	log.Println("Raw text stored in S3:", bucketName, "/", event.OutputFileName)
	return "Success", nil
}

func main() {
	log.Println("Starting handler...")
	lambda.Start(handler)
}
