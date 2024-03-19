build:
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/process-documents/bootstrap -ldflags "-s -w" ./process-documents
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/convert-files/bootstrap -ldflags "-s -w" ./convert-files
	GOARCH=arm64 GOOS=linux go build -o ../out/lambda/extract-english/bootstrap -ldflags "-s -w" ./extract-english
bd:
	make build
	cd ../infra && cdk deploy
	cd ../functions
