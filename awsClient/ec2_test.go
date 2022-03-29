package awsClient

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/camaeell/aws-asg-dyndns/awsClient/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestPrivateGetInstanceIps(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{
			{
				PrivateIpAddress: &expectedPrivateIp,
				Association: &types.NetworkInterfaceAssociation{
					PublicIp: &expectedPublicIp,
				},
			},
		},
	}

	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, expectedErr)

	instanceId := "i-fake123"

	privateIp, publicIp, err := getInstanceIps(ctx, m, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestPrivateGetInstanceIpsError(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	instanceId := "i-fake123"
	expectedErr := fmt.Errorf("Can't obtain ENI details for instance %s, Empty Interfaces list for instance: %s", instanceId, instanceId)
	returnedErr := fmt.Errorf("Empty Interfaces list for instance: %s", instanceId)

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{},
	}
	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, returnedErr)

	var privateIp, publicIp *string

	assert.PanicsWithValuef(t, expectedErr.Error(), func() { privateIp, publicIp, _ = getInstanceIps(ctx, m, instanceId) }, "The code should panic with %s", expectedErr)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestPrivateGetInstanceIpsFromTags(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil
	instanceId := "i-fake123"

	privateIpToken := "privateIp"
	publicIpToken := "publicIp"

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseGetTags := ec2.DescribeTagsOutput{
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
	}

	m.EXPECT().DescribeTags(gomock.Any(), gomock.Any()).
		Return(&fakeResponseGetTags, expectedErr)

	privateIp, publicIp, err := getInstanceIpsFromTags(ctx, m, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestPrivateGetInstanceIpsFromTagsError(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	var expectedErr error = fmt.Errorf("AWS dummy error")
	instanceId := "i-fake123"

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseGetTags := ec2.DescribeTagsOutput{
		Tags: []types.TagDescription{},
	}
	m.EXPECT().DescribeTags(gomock.Any(), gomock.Any()).
		Return(&fakeResponseGetTags, expectedErr)

	privateIp, publicIp, err := getInstanceIpsFromTags(ctx, m, instanceId)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestGetInstanceIpsPositive(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{
			{
				PrivateIpAddress: &expectedPrivateIp,
				Association: &types.NetworkInterfaceAssociation{
					PublicIp: &expectedPublicIp,
				},
			},
		},
	}

	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, expectedErr)

	instanceId := "i-fake123"

	privateIp, publicIp, err := GetInstanceIps(ctx, m, instanceId, false)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestGetInstanceIpsErr(t *testing.T) {
	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	instanceId := "i-fake123"
	expectedErr := fmt.Errorf("Can't obtain ENI details for instance %s, Empty Interfaces list for instance: %s", instanceId, instanceId)
	returnedErr := fmt.Errorf("Empty Interfaces list for instance: %s", instanceId)

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{},
	}

	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, returnedErr)

	var privateIp, publicIp *string

	assert.PanicsWithValuef(t, expectedErr.Error(), func() { privateIp, publicIp, _ = GetInstanceIps(ctx, m, instanceId, false) }, "The code should panic with %s", expectedErr)

	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}

func TestGetInstanceIpsFromTagsPositive(t *testing.T) {

	expectedPrivateIp := "192.168.199.99"
	expectedPublicIp := "99.100.101.102"
	var expectedErr error = nil
	privateIpToken := "privateIp"
	publicIpToken := "publicIp"
	instanceId := "i-fake123"

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{},
	}
	fakeResponseGetTags := ec2.DescribeTagsOutput{
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
	}
	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, expectedErr)
	m.EXPECT().DescribeTags(gomock.Any(), gomock.Any()).
		Return(&fakeResponseGetTags, expectedErr)

	privateIp, publicIp, err := GetInstanceIps(ctx, m, instanceId, true)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, *privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, *privateIp)
	assert.Equalf(t, expectedPublicIp, *publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, *publicIp)
}

func TestGetInstanceIpsFromTagsNegative(t *testing.T) {

	var expectedPrivateIp *string = nil
	var expectedPublicIp *string = nil
	var expectedErr error = nil
	instanceId := "i-fake123"

	ctrl := gomock.NewController(t)
	m := mocks.NewMockEC2API(ctrl)
	ctx := context.TODO()

	fakeResponseNI := ec2.DescribeNetworkInterfacesOutput{
		NetworkInterfaces: []types.NetworkInterface{},
	}
	fakeResponseGetTags := ec2.DescribeTagsOutput{
		Tags: []types.TagDescription{},
	}

	m.EXPECT().DescribeNetworkInterfaces(gomock.Any(), gomock.Any()).
		Return(&fakeResponseNI, expectedErr)
	m.EXPECT().DescribeTags(gomock.Any(), gomock.Any()).
		Return(&fakeResponseGetTags, expectedErr)

	privateIp, publicIp, err := GetInstanceIps(ctx, m, instanceId, true)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedPrivateIp, privateIp, "Wrong privateIp. Expected %s, got %s", expectedPrivateIp, privateIp)
	assert.Equalf(t, expectedPublicIp, publicIp, "Wrong publicIp. Expected %s, got %s", expectedPublicIp, publicIp)
}
