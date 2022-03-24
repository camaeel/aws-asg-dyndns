package awsClient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Ec2FakeClient struct {
	DescribeTagsMock *ec2.DescribeTagsOutput
	CreateTagsMock   *ec2.CreateTagsOutput
	DescribeNIMock   *ec2.DescribeNetworkInterfacesOutput
}

func (c *Ec2FakeClient) DescribeTags(ctx context.Context, params *ec2.DescribeTagsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeTagsOutput, error) {
	return c.DescribeTagsMock, nil
}

func (c *Ec2FakeClient) CreateTags(ctx context.Context, params *ec2.CreateTagsInput, optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error) {
	return c.CreateTagsMock, nil
}
func (c *Ec2FakeClient) DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	return c.DescribeNIMock, nil
}
