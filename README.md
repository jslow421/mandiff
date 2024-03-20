# mandiff

To use the system as it is currently deployed follow these instructions (will need to update when system is moved)

## Instructions for textract system use

- Be authed for CLI access
- Place documents in bucket/folder (med-manual-upload-bucket for now)
- Update the folder name here

```bash
aws lambda invoke --function-name SlowikDocumentStack-ProcessDocuments2119A556-wkw9GtCKs1rg --payload '{"objectPrefix":"FOLDER_HERE/"}' --cli-binary-format raw-in-base64-out /dev/stdout
```

- You'll receive textract job IDs

- Check to see when they're done ex:

```bash
aws textract get-document-analysis --job-id "YOUR_ID_HERE"
```

- Once the job is complete copy the first job ID and insert it here

```bash
aws lambda invoke --function-name SlowikDocumentStack-CreateTextFileD8BBA755-6IRP9v6jrCVj --payload '{"jobId":"JOB_ID_HERE", "outputFileName": "FILE_NAME_HERE", "shouldFilterLanguage": false}' --cli-binary-format raw-in-base64-out /dev/stdout
```

- This will generate the flat text file
- Find the output with your name in the bucket `med-manual-complete-bucket`

- run LLM function

```bash
aws lambda invoke --function-name SlowikDocumentStack-LlmFunction49CBF322-BYpffuwiGVR5 --payload '{"documentBucket":"med-manual-complete-bucket"}' --cli-binary-format raw-in-base64-out /dev/stdout
```
