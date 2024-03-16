package main

import (
	"fmt"
	"log"

	runtime "github.com/aws/aws-lambda-go/lambda"
)

func handler() {
	fmt.Println("Processing documents...")
}

func main() {
	log.Println("Starting handler...")
	runtime.Start(handler)
}
