package awsClient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func GetSSMParameterValue(ctx context.Context, ssmClient *ssm.Client, parameterName string) (string, error) {
	input := ssm.GetParametersInput{
		WithDecryption: true,
		Names:          []string{},
	}

	output, err := ssmClient.GetParameters(ctx, &input)
	if err != nil {
		return "", err
	}
	ret := *output.Parameters[0].Value

	return ret, err
}
