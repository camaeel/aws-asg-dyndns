package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

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
