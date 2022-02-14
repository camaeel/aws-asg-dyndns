package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	for i := range event.Records {
		err = processRecord(ctx, client, event.Records[i])
		if err != nil {
			return "Fail", err
		}
	}

	return "Ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}

func processRecord(ctx context.Context, client *ec2.Client, record events.SQSMessage) error {
	var result error

	recordLog, err := json.Marshal(record)
	if err != nil {
		log.Fatal("Can't marshal event")
		return err
	}
	log.Printf("Record: %s", recordLog)

	var bodyJson LifecycleMessage
	err = json.Unmarshal([]byte(record.Body), &bodyJson)
	if err != nil {
		log.Fatal("Can't unmarshal event body", err)
		return err
	}
	log.Printf("Record body: %s", bodyJson)

	return result
}

func getInstanceIps(ctx context.Context, client *ec2.Client, instanceId string) (*string, *string, error) {

	filterName := "attachment.instance-id"

	var filter types.Filter = types.Filter{Name: &filterName, Values: []string{instanceId}}
	var nIInput ec2.DescribeNetworkInterfacesInput = ec2.DescribeNetworkInterfacesInput{Filters: []types.Filter{filter}}
	nIOutput, err := client.DescribeNetworkInterfaces(ctx, &nIInput)
	if err != nil {
		panic("Can't obtain ENI details for instance " + instanceId + ", " + err.Error())
	}

	if len(nIOutput.NetworkInterfaces) < 1 {
		return nil, nil, errors.New("Empty Interfaces list for instance: " + instanceId)
	}

	if len(nIOutput.NetworkInterfaces) > 1 {
		log.Printf("Warning! Instance: %s has %d network interfaces. Using first.", instanceId, len(nIOutput.NetworkInterfaces))
	}

	privateIp := nIOutput.NetworkInterfaces[0].PrivateIpAddress

	var publicIp *string = nil
	if nIOutput.NetworkInterfaces[0].Association != nil && nIOutput.NetworkInterfaces[0].Association.PublicIp != nil {
		publicIp = nIOutput.NetworkInterfaces[0].Association.PublicIp
	}

	return privateIp, publicIp, nil
}
