package awsClient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SsmFakeClient struct {
	GetParametersMock *ssm.GetParameterOutput
	GetParametersErr  error
}

func (c *SsmFakeClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return c.GetParametersMock, c.GetParametersErr
}
