package cloudflare_dns

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	awsClient_mocks "github.com/camaeell/aws-asg-dyndns/awsClient/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCloudflareDomainToZoneResolve(t *testing.T) {
	domain := "test.test1.example.com"
	expected := "example.com"

	var cfp cloudflareProvider
	cfp.detectZone(domain)

	assert.Equal(t, expected, cfp.zone, "Zone should be equal.")
}

func TestNewCloudflareProvider(t *testing.T) {
	ctx := context.TODO()
	expectedToken := "token123ABC"
	expectedName := "TOKEN_PARAM_NAME"
	domain := "www.test.example.com"
	expectedZone := "example.com"

	ctrl := gomock.NewController(t)
	m := awsClient_mocks.NewMockSSMAPI(ctrl)

	fakeSSMResponse := ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			Name:  &expectedName,
			Value: &expectedToken,
		},
	}
	m.EXPECT().GetParameter(gomock.Any(), gomock.Any()).Return(&fakeSSMResponse, nil)

	provider, err := NewCloudflareProvider(ctx, m, domain)

	assert.Equalf(t, expectedToken, provider.token, "Token %s not equal to expected value %s", provider.token, expectedToken)
	assert.Equalf(t, expectedZone, provider.zone, "Zone %s not equal to expected value %s", provider.zone, expectedZone)
	assert.Nilf(t, err, "Error %s not null", err)
}
