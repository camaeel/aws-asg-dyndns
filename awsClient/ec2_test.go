package awsClient

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestPrivateGetInstanceIpsTest(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{
				{
					PrivateIpAddress: &expectedPrivateIp,
					Association: &types.NetworkInterfaceAssociation{
						PublicIp: &expectedPublicIp,
					},
				},
			},
		},
	}
	ctx := context.TODO()
	instanceId := "i-fake123"

	privateIp, publicIp, err := getInstanceIps(ctx, &fakeClient, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestPrivateGetInstanceIpsErrorTest(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	instanceId := "i-fake123"
	var expectedErr error = fmt.Errorf("Empty Interfaces list for instance: %s", instanceId)

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{},
		},
	}
	ctx := context.TODO()

	privateIp, publicIp, err := getInstanceIps(ctx, &fakeClient, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestPrivateGetInstanceIpsFromTagsTest(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil
	instanceId := "i-fake123"

	privateIpToken := "privateIp"
	publicIpToken := "publicIp"

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeTagsMock: &ec2.DescribeTagsOutput{
			Tags: []types.TagDescription{
				{
					Key:        &privateIpToken,
					Value:      &expectedPrivateIp,
					ResourceId: &instanceId,
				},
				{
					Key:        &publicIpToken,
					Value:      &expectedPublicIp,
					ResourceId: &instanceId,
				},
			},
		},
	}
	ctx := context.TODO()

	privateIp, publicIp, err := getInstanceIpsFromTags(ctx, &fakeClient, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestPrivateGetInstanceIpsFromTagsErrorTest(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	var expectedErr error = fmt.Errorf("AWS dummy error")
	instanceId := "i-fake123"

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeTagsMockErr: expectedErr,
	}
	ctx := context.TODO()

	privateIp, publicIp, err := getInstanceIpsFromTags(ctx, &fakeClient, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestGetInstanceIpsPositiveTest(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{
				{
					PrivateIpAddress: &expectedPrivateIp,
					Association: &types.NetworkInterfaceAssociation{
						PublicIp: &expectedPublicIp,
					},
				},
			},
		},
	}
	ctx := context.TODO()
	instanceId := "i-fake123"

	privateIp, publicIp, err := GetInstanceIps(ctx, &fakeClient, instanceId, false)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestGetInstanceIpsErrTest(t *testing.T) {
	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	instanceId := "i-fake123"
	var expectedErr error = fmt.Errorf("Empty Interfaces list for instance: %s", instanceId)

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{},
		},
	}
	ctx := context.TODO()

	privateIp, publicIp, err := GetInstanceIps(ctx, &fakeClient, instanceId, false)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestGetInstanceIpsFromTagsPositiveTest(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil
	privateIpToken := "privateIp"
	publicIpToken := "publicIp"
	instanceId := "i-fake123"

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{},
		},
		DescribeTagsMock: &ec2.DescribeTagsOutput{
			Tags: []types.TagDescription{
				{
					Key:        &privateIpToken,
					Value:      &expectedPrivateIp,
					ResourceId: &instanceId,
				},
				{
					Key:        &publicIpToken,
					Value:      &expectedPublicIp,
					ResourceId: &instanceId,
				},
			},
		},
	}
	ctx := context.TODO()

	privateIp, publicIp, err := GetInstanceIps(ctx, &fakeClient, instanceId, true)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestGetInstanceIpsFromTagsNegativeTest(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	var expectedErr error = nil
	instanceId := "i-fake123"

	var fakeClient Ec2FakeClient = Ec2FakeClient{
		DescribeNIMock: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []types.NetworkInterface{},
		},
		DescribeTagsMock: &ec2.DescribeTagsOutput{
			Tags: []types.TagDescription{},
		},
	}
	ctx := context.TODO()

	privateIp, publicIp, err := GetInstanceIps(ctx, &fakeClient, instanceId, true)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}
