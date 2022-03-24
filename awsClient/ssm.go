package awsClient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetSSMParameterValue(ctx context.Context, ssmClient SSMAPI, parameterName string) (string, error) {
	input := ssm.GetParameterInput{
		WithDecryption: true,
		Name:           &parameterName,
	}

	output, err := ssmClient.GetParameter(ctx, &input)
	if err != nil {
		return "", err
	}
	ret := *output.Parameter.Value

	return ret, err
}
