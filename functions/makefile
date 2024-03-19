build:
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/process-documents/bootstrap -ldflags "-s -w" ./process-documents
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/convert-files/bootstrap -ldflags "-s -w" ./convert-files
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/llm/bootstrap -ldflags "-s -w" ./llm
bd:
	make build && cd ../infra && cdk deploy && cd ../functions
	