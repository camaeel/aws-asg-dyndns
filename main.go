package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	ec2Client := ec2.NewFromConfig(cfg)
	autoscalingClient := autoscaling.NewFromConfig(cfg)

	for i := range event.Records {
		err = processRecord(ctx, ec2Client, autoscalingClient, event.Records[i])
		if err != nil {
			return "Fail", err
		}
	}

	return "Ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}

func processRecord(ctx context.Context, ec2Client *ec2.Client, autoscalingClient *autoscaling.Client, record events.SQSMessage) error {
	var privateIp, publicIp *string

	recordLog, err := json.Marshal(record)
	if err != nil {
		log.Fatal("Error! Can't marshal event")
		return err
	}
	log.Printf("Record: %s", recordLog)

	var recordBody LifecycleMessage
	err = json.Unmarshal([]byte(record.Body), &recordBody)
	if err != nil {
		log.Fatal("Error! Can't unmarshal event body", err)
		return err
	}

	log.Printf("Record body: %s", recordBody)

	if recordBody.LifecycleTransition == "autoscaling:EC2_INSTANCE_LAUNCHING" {
		if recordBody.EC2InstanceId == "" {
			return errors.New("body.EC2InstanceId not set")
		}

		privateIp, publicIp, err = getInstanceIps(ctx, ec2Client, recordBody.EC2InstanceId)
		if err != nil {
			log.Fatal("Error! Can't obtain IPs", err)
			return err
		}

		log.Printf("Instance: %s, got IPs: privateIp: %s, publicIp: %s", recordBody.EC2InstanceId, *privateIp, *publicIp)

		err = tagResource(ctx, ec2Client, recordBody.EC2InstanceId, privateIp, publicIp)
		if err != nil {
			log.Fatal("Error! Can't tag IPs on instance", err)
			return err
		}
		log.Printf("Instance %s tagged with IPs", recordBody.EC2InstanceId)

		// TODO: handle cloudflare add
		err = completeLifecycleHook(ctx, autoscalingClient, recordBody)
		if err != nil {
			log.Fatal("Error! Can't complete lifecyle hook. ", err)
			return err
		} else {
			log.Printf("Lifecycle hook completed successful.")
		}

	} else if recordBody.LifecycleTransition == "autoscaling:EC2_INSTANCE_TERMINATING" {
		privateIp, publicIp, err = getInstanceIps(ctx, ec2Client, recordBody.EC2InstanceId)
		if err != nil {
			log.Print("Warning! Can't obtain IPs from instance. Will try with tags. Err:", err)
			privateIp, publicIp, err = getInstanceIpsFromTags(ctx, ec2Client, recordBody.EC2InstanceId)
			if err != nil {
				log.Fatal("Error! Can't get IPs from instance's tags. Err:", err)
				return err
			}
		}
		log.Printf("Instance: %s, got IPs: privateIp: %s, publicIp: %s", recordBody.EC2InstanceId, *privateIp, *publicIp)

		// TODO: handle cloudflare remove

		err = completeLifecycleHook(ctx, autoscalingClient, recordBody)
		if err != nil {
			log.Print("Error! Can't complete lifecyle hook. ", err)
			return err
		} else {
			log.Printf("Lifecycle hook completed successful.")
		}
	} else {
		log.Printf("Warning! Unknown LifecycleTransition: [%s]", recordBody.LifecycleTransition)
		return nil
	}

	return nil
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

func getInstanceIpsFromTags(ctx context.Context, client *ec2.Client, instanceId string) (*string, *string, error) {
	var RESOURCE_ID string = "resource-id"
	var KEY_STR string = "key"
	var privateIp, publicIp *string

	var input ec2.DescribeTagsInput = ec2.DescribeTagsInput{
		Filters: []types.Filter{
			{Name: &RESOURCE_ID, Values: []string{instanceId}},
			{Name: &KEY_STR, Values: []string{"privateIp", "publicIp"}},
		},
	}

	output, err := client.DescribeTags(ctx, &input)

	if err != nil {
		log.Printf("ERROR: Can't obtain IPs from tags for instance: %s. Err: %s", instanceId, err)
		return nil, nil, err
	}

	for i := range output.Tags {
		if *output.Tags[i].Key == "privateIp" {
			privateIp = output.Tags[i].Value
		} else if *output.Tags[i].Key == "publicIp" {
			publicIp = output.Tags[i].Value
		} else {
			log.Printf("WARN: Unexpected tag key: %s", *output.Tags[i].Key)
		}
	}

	return privateIp, publicIp, nil
}

func completeLifecycleHook(ctx context.Context, client *autoscaling.Client, message LifecycleMessage) error {
	var actionResult string = "CONTINUE"
	var input autoscaling.CompleteLifecycleActionInput = autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  &message.AutoScalingGroupName,
		LifecycleHookName:     &message.LifecycleHookName,
		LifecycleActionResult: &actionResult,
		LifecycleActionToken:  &message.LifecycleActionToken,
	}

	_, err := client.CompleteLifecycleAction(ctx, &input)

	return err
}

func tagResource(ctx context.Context, client *ec2.Client, instanceId string, privateIp *string, publicIp *string) error {
	var privateIpKey string = "privateIp"
	var publicIpKey string = "publicIp"
	var input ec2.CreateTagsInput = ec2.CreateTagsInput{
		Resources: []string{instanceId},
		Tags: []types.Tag{
			{Key: &privateIpKey, Value: privateIp},
			{Key: &publicIpKey, Value: publicIp},
		},
	}

	_, err := client.CreateTags(ctx, &input)

	return err
}
