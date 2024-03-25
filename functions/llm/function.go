package main

import (
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/bedrock"
)

type LlmEvent struct {
	JobId           string `json:"jobId"`
	InitialDocument string `json:"initialDocument"`
	UpdatedDocument string `json:"updatedDocument"`
	DocumentBucket  string `json:"documentBucket"`
	Document1Key    string `json:"document1Key"`
	Document2Key    string `json:"document2Key"`
	CvText          string `json:"cvText"`
}

const MAX_HISTORY_LENGTH = 10
const DOCUMENT_COMPARE_TEMPLATE = `
You are an assistant tasked with comparing two documents and noting the differences. These documents are manuals from medical equipment.
They will be parsed with textract to provide only the raw text values.
If languages other than English are present, please only evaluate and output items in English.
If there are other languages present, please ignore them.

The individuals you are preparing the comparison for are all medical professionals, 
and they require knowledge of the specific differences between the two so that they can effectively do their jobs.
This includes, but is not necessarily limited to, information about use, cleaning, care, and maintenance of the equipment.

Specific areas of interest include anything related to:
document ID and publish date
cleaning/sterilization
preventative maintenance
troubleshooting
testing
wiring
schematics
parts list
calibration

If the changes are only related to language, punctuation, or formatting, it is ok to simply state that, 
but if there are fundamental differences that a medical professional would want to know about please explicitly state them.

I will provide the documents separately. They will be inside xml tags called <document1> and <document2>.
If you can discern which document is more recently created, please refer to that as "the updated" document.


<document1>
{{ .Document1 }}
</document1>


<document2>
{{ .Document2 }}
</document2>

History: {{ .History }}
Additional input: {{ .Input }}

Please confirm your understanding of the above requirements

Assistant: I understand the requirements.
I will read the english manuals.
Differences that I should point out are related to information about use, cleaning, care, and maintenance of the equipment mentioned in the manual.
I will inform of any differences I find, and I will explain what changed and where.
If I can determine which document is more recent, I will refer to that as "the updated" document.

Human: Read back is correct, please tell me what differences you see.
Assistant:
`

const EDUCATION_EXTRACTION_TEMPLATE = `
You are an assistant tasked with comparing two documents and noting the differences. These documents are manuals from medical equipment.
They will be parsed with textract to provide only the raw text values.
If languages other than English are present, please only evaluate and output items in English.
If there are other languages present, please ignore them.

The individuals you are preparing the comparison for are all medical professionals, 
and they require knowledge of the specific differences between the two so that they can effectively do their jobs.
This includes, but is not necessarily limited to, information about use, cleaning, care, and maintenance of the equipment.

Specific areas of interest include anything related to:
document ID and publish date
cleaning/sterilization
preventative maintenance
troubleshooting
testing
wiring
schematics
parts list
calibration

If the changes are only related to language, punctuation, or formatting, it is ok to simply state that, 
but if there are fundamental differences that a medical professional would want to know about please explicitly state them.

I will provide the documents separately. They will be inside xml tags called <document1> and <document2>.
If you can discern which document is more recently created, please refer to that as "the updated" document.


<document1>
{{ .CVText }}
</document1>

History: {{ .History }}
Additional input:

Please confirm your understanding of the above requirements

Assistant: I understand the requirements.
I will read the english manuals.
Differences that I should point out are related to information about use, cleaning, care, and maintenance of the equipment mentioned in the manual.
I will inform of any differences I find, and I will explain what changed and where.
If I can determine which document is more recent, I will refer to that as "the updated" document.

Human: Read back is correct, please tell me what differences you see.
Assistant:
`

func buildDocumentComparisonPrompt(document1 string, document2 string) (string, error) {
	template, err := template.New("prompt").Parse(DOCUMENT_COMPARE_TEMPLATE)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := struct {
		History   string
		Input     string
		Document1 string
		Document2 string
	}{
		History:   "",
		Input:     "",
		Document1: document1,
		Document2: document2,
	}

	var builder strings.Builder

	err = template.Execute(&builder, data)

	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return builder.String(), nil
}

func buildEducationExtractionPrompt() (string, error) {
	template, err := template.New("prompt").Parse(EDUCATION_EXTRACTION_TEMPLATE)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	data := struct {
		History string
		CVText  string
	}{
		History: "",
		CVText:  "",
	}

	var builder strings.Builder

	err = template.Execute(&builder, data)

	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return builder.String(), nil
}

func getTextFromS3File(ctx context.Context, bucket string, key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load default config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	getObject := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	result, err := s3Client.GetObject(ctx, getObject)
	if err != nil {
		log.Fatalf("Failed to get object: %v", err)
		return "", err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		log.Fatalf("Failed to read object body: %v", err)
		return "", err
	}

	return string(data), nil
}

// Ok, enums might be nice :)
type PromptType int

const (
	DocumentComparison PromptType = iota
	EducationExtraction
)

func (p PromptType) String() string {
	return [...]string{"Document Comparison", "Education Extraction"}[p]
}

type PromptResult struct {
	Prompt     string
	Error      error
	PromptType PromptType
}

// Build a prompt based on the input data
func promptFactory(firstDocumentText string, secondDocumentText string, cvText string) PromptResult {
	if cvText != "" {
		prompt, promptErr := buildEducationExtractionPrompt()
		return PromptResult{Prompt: prompt, Error: promptErr, PromptType: EducationExtraction}
	} else {
		prompt, promptErr := buildDocumentComparisonPrompt(firstDocumentText, secondDocumentText)
		return PromptResult{Prompt: prompt, Error: promptErr, PromptType: DocumentComparison}
	}
}

func handler(ctx context.Context, event *LlmEvent) (string, error) {
	var firstDocumentText string
	var SecondDocumentText string
	bucket := event.DocumentBucket
	llm, err := bedrock.New(
		bedrock.WithModel("anthropic.claude-3-sonnet-20240229-v1:0"),
	)

	if err != nil {
		log.Fatalf("Failed to create LLM: %v", err)
		return "", err
	}

	var wg = sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		firstDocumentText, _ = getTextFromS3File(ctx, bucket, event.Document1Key)
	}()

	go func() {
		defer wg.Done()
		SecondDocumentText, _ = getTextFromS3File(ctx, bucket, event.Document2Key)
	}()

	wg.Wait()

	prompt := promptFactory(firstDocumentText, SecondDocumentText, event.CvText)

	if prompt.Error != nil {
		log.Fatalf("Failed to build prompt: %v", prompt.Error)
		return "", prompt.Error
	}

	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt.Prompt,
		llms.WithMaxTokens(180000),
		llms.WithMaxLength(50000),
	)

	if err != nil {
		log.Fatalf("Failed to generate completion: %v", err)
		return "", err
	}

	return completion, nil
}

func main() {
	log.Println("Starting LLM...")
	lambda.Start(handler)
	log.Println("LLM completed.")
}
