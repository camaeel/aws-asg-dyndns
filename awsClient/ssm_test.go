package awsClient

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
)

func TestGetSSMParameterValue(t *testing.T) {

	expectedParamValue := "secretToken123"
	parameterName := "testing-param1"
	var expectedErr error = nil

	var fakeClient SsmFakeClient = SsmFakeClient{
		GetParametersMock: &ssm.GetParameterOutput{
			Parameter: &types.Parameter{
				Name:  &parameterName,
				Value: &expectedParamValue,
			},
		},
	}
	ctx := context.TODO()
	paramValue, err := GetSSMParameterValue(ctx, &fakeClient, parameterName)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedParamValue, paramValue, "Wrong paramValue. Expected %s, got %s", expectedParamValue, paramValue)
}

func TestGetSSMParameterValueError(t *testing.T) {

	expectedParamValue := ""
	parameterName := "testing-param1"
	var expectedErr error = fmt.Errorf("Can't find data")

	var fakeClient SsmFakeClient = SsmFakeClient{
		GetParametersMock: &ssm.GetParameterOutput{},
		GetParametersErr:  expectedErr,
	}
	ctx := context.TODO()
	paramValue, err := GetSSMParameterValue(ctx, &fakeClient, parameterName)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedParamValue, paramValue, "Wrong paramValue. Expected %s, got %s", expectedParamValue, paramValue)
}
