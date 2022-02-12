package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	for i := range event.Records {
		processRecord(event.Records[i])
	}

	return "Ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}

func processRecord(record events.SQSMessage) error {
	var result error

	recordLog, err := json.Marshal(record)
	if err != nil {
		log.Fatal("Can't marshal event")
		return err
	}
	log.Printf("Record: %s", recordLog)

	return result
}
