package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
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

func retrieveDocumentsFomS3(event DocumentEvent) {
	log.Println("Retrieving documents from S3...")

	cfg, _ := config.LoadDefaultConfig(context.TODO())

	s3Client := s3.NewFromConfig(cfg)
	textractClient := textract.NewFromConfig(cfg)
	uploadBucket := os.Getenv("UPLOAD_BUCKET")
	delimiter := "/"

	bucketObjects, _ := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    &uploadBucket,
		Prefix:    &event.ObjectPrefix,
		Delimiter: &delimiter,
	})

	log.Println("Found documents in bucket: ", len(bucketObjects.Contents))
	for _, obj := range bucketObjects.Contents {
		log.Printf("Found object with key: %s", *obj.Key)
	}

	if len(bucketObjects.Contents) > 3 {
		log.Fatalf("Too many documents in bucket: %d", len(bucketObjects.Contents))
		return
	}

	for _, obj := range bucketObjects.Contents {
		log.Printf("Processing document: %s", *obj.Key)
	}

	for i, document := range bucketObjects.Contents {
		log.Printf("Processing document: %s", *document.Key)

		if !strings.Contains(*document.Key, ".") {
			log.Println("Skipping 'folder' key: ", *document.Key)
		}

		objKey := *bucketObjects.Contents[i].Key
		objBucket := *bucketObjects.Name
		objOutputKey := objKey + "-complete"
		snsTopicArn := os.Getenv("SNS_TOPIC_ARN")
		roleArn := os.Getenv("ROLE_ARN")
		outputBucket := os.Getenv("COMPLETE_BUCKET")

		input := &textract.StartDocumentAnalysisInput{
			DocumentLocation: &types.DocumentLocation{
				S3Object: &types.S3Object{
					Bucket: &objBucket,
					Name:   &objKey,
				},
			},
			OutputConfig: &types.OutputConfig{
				S3Bucket: &outputBucket,
				S3Prefix: &objOutputKey,
			},
			NotificationChannel: &types.NotificationChannel{
				SNSTopicArn: &snsTopicArn,
				RoleArn:     &roleArn,
			},
			FeatureTypes: []types.FeatureType{
				types.FeatureTypeTables,
			},
		}

		job, err := textractClient.StartDocumentAnalysis(context.TODO(), input)

		if err != nil {
			log.Fatalf("Error processing document %s: %v", *document.Key, err)
		}

		log.Printf("Started processing document %s with job id: %s", *document.Key, *job.JobId)
	}
}

type DocumentEvent struct {
	JobId            string `json:"jobId"`
	InputBucketName  string `json:"inputBucketName"`
	OutputBucketName string `json:"outputBucketName"`
	OutputFileName   string `json:"outputFileName"`
	ObjectPrefix     string `json:"objectPrefix"`
}

func handler(ctx context.Context, event *DocumentEvent) {
	fmt.Println("Processing documents...")
	retrieveDocumentsFomS3(*event)
	log.Println("Documents sent for processing")
}

func main() {
	log.Println("Starting handler...")
	lambda.Start(handler)
}
