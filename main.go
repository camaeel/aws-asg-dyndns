package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	eventJson, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		log.Fatal("Can't log event")
		return "", err
	}

	log.Printf("EVENT: %s", eventJson)
	// environment variables
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	// log.Println("ALL ENV VARS:")
	// for _, element := range os.Environ() {
	// 	log.Println(element)
	// }

	return "Ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
