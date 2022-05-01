package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/camaeell/aws-asg-dyndns/awsClient"
	"github.com/camaeell/aws-asg-dyndns/dns"
)

func HandleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	ec2Client := ec2.NewFromConfig(cfg)
	autoscalingClient := autoscaling.NewFromConfig(cfg)
	ssmClient := ssm.NewFromConfig(cfg)

	for i := range event.Records {
		err = processRecord(ctx, ec2Client, autoscalingClient, ssmClient, event.Records[i])
		if err != nil {
			return "Fail", err
		}
	}

	return "Ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}

func processRecord(ctx context.Context, ec2Client *ec2.Client, autoscalingClient *autoscaling.Client, ssmClient *ssm.Client, record events.SQSMessage) error {
	recordLog, err := json.Marshal(record)
	if err != nil {
		log.Fatal("Error! Can't marshal event")
		return err
	}
	log.Printf("Record: %s", recordLog)

	var recordBody awsClient.LifecycleMessage
	err = json.Unmarshal([]byte(record.Body), &recordBody)
	if err != nil {
		log.Fatal("Error! Can't unmarshal event body", err)
		return err
	}

	log.Printf("Record body: %s", recordBody)

	if recordBody.LifecycleTransition == "autoscaling:EC2_INSTANCE_LAUNCHING" {
		err = processInstanceLaunchRecord(ctx, ec2Client, ssmClient, recordBody)
		if err != nil {
			return err
		}
	} else if recordBody.LifecycleTransition == "autoscaling:EC2_INSTANCE_TERMINATING" {
		err = processInstanceTerminateRecord(ctx, ec2Client, ssmClient, recordBody)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Error! Unknown LifecycleTransition: [%s]", recordBody.LifecycleTransition)
	}

	err = awsClient.CompleteLifecycleHook(ctx, autoscalingClient, recordBody)
	if err != nil {
		log.Print("Error! Can't complete lifecyle hook. ", err)
		return err
	} else {
		log.Printf("Lifecycle hook completed successful.")
	}

	return nil
}

func processInstanceLaunchRecord(ctx context.Context, ec2Client *ec2.Client, ssmClient *ssm.Client, recordBody awsClient.LifecycleMessage) error {
	var privateIp, publicIp *string

	if recordBody.EC2InstanceId == "" {
		return errors.New("body.EC2InstanceId not set")
	}

	privateIp, publicIp, err := awsClient.GetInstanceIps(ctx, ec2Client, recordBody.EC2InstanceId, false)
	if err != nil {
		log.Fatal("Error! Can't obtain IPs", err)
		return err
	}

	log.Printf("Instance added: %s, got IPs: privateIp: %s, publicIp: %s", recordBody.EC2InstanceId, *privateIp, *publicIp)

	err = awsClient.TagEC2Instance(ctx, ec2Client, recordBody.EC2InstanceId, privateIp, publicIp)
	if err != nil {
		log.Fatal("Error! Can't tag IPs on instance", err)
		return err
	}
	log.Printf("Instance %s tagged with IPs", recordBody.EC2InstanceId)

	if recordBody.NotificationMetadata != nil && recordBody.NotificationMetadata["domainList"] != nil {
		for i := range recordBody.NotificationMetadata["domainList"] {
			err := dns.DnsEntryAddIp(ctx, ssmClient, recordBody.NotificationMetadata["domainList"][i], publicIp)
			if err != nil {
				log.Printf("Error! Can't add DNS entry for ip: %s, domain: %s. %s", *publicIp, recordBody.NotificationMetadata["domainList"][i], err)
				return err
			}
		}
	}

	return nil
}

func processInstanceTerminateRecord(ctx context.Context, ec2Client *ec2.Client, ssmClient *ssm.Client, recordBody awsClient.LifecycleMessage) error {
	var privateIp, publicIp *string

	if recordBody.EC2InstanceId == "" {
		return errors.New("body.EC2InstanceId not set")
	}

	privateIp, publicIp, err := awsClient.GetInstanceIps(ctx, ec2Client, recordBody.EC2InstanceId, true)
	if err != nil {
		log.Fatal("Error! Can't tag IPs on instance", err)
		return err
	}

	log.Printf("Instance removed: %s, got IPs: privateIp: %s, publicIp: %s", recordBody.EC2InstanceId, *privateIp, *publicIp)

	if recordBody.NotificationMetadata != nil && recordBody.NotificationMetadata["domainList"] != nil {
		for i := range recordBody.NotificationMetadata["domainList"] {
			err := dns.DnsEntryRemoveIp(ctx, ssmClient, recordBody.NotificationMetadata["domainList"][i], publicIp)
			if err != nil {
				log.Printf("Error! Can't remove DNS entry for ip: %s, domain: %s. %s", *publicIp, recordBody.NotificationMetadata["domainList"][i], err)
				return err
			}
		}
	}

	return nil
}
