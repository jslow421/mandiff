package main

import (
	"context"
	"fmt"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type BucketBasics struct {
	S3Client *s3.Client
}

type Document struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func retrieveDocumentsFomS3() {
	log.Println("Retrieving documents from S3...")
	cfg, _ := config.LoadDefaultConfig(context.TODO())

	//s3Client := s3.NewFromConfig(cfg)
	textractClient := textract.NewFromConfig(cfg)

	document1 := Document{
		Bucket: "med-manual-upload-bucket",
		Key:    "Manual-2016.pdf",
	}

	document2 := Document{
		Bucket: "med-manual-upload-bucket",
		Key:    "Manual-2020.pdf",
	}

	documents := []Document{document1, document2}

	for _, document := range documents {
		log.Printf("Processing document: %s", document.Key)
		input := &textract.StartDocumentAnalysisInput{
			DocumentLocation: &types.DocumentLocation{
				S3Object: &types.S3Object{
					Bucket: &document.Bucket,
					Name:   &document.Key,
				},
			},
			OutputConfig: &types.OutputConfig{
				S3Bucket: aws.String("med-manual-complete-bucket"),
				S3Prefix: aws.String(document.Key + "-complete"),
			},
			FeatureTypes: []types.FeatureType{"TABLES", "FORMS"},
		}

		job, err := textractClient.StartDocumentAnalysis(context.TODO(), input)

		if err != nil {
			log.Fatalf("Error processing document %s: %v", document.Key, err)
		}

		log.Printf("Started processing document %s with job id: %s", document.Key, *job.JobId)
	}
}

func handler() {
	fmt.Println("Processing documents...")
	retrieveDocumentsFomS3()
	log.Println("Documents sent for processing")
}

func main() {
	log.Println("Starting handler...")
	runtime.Start(handler)
}
