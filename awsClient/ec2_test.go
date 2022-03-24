package awsClient

import (
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
