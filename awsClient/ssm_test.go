package awsClient

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/camaeell/aws-asg-dyndns/awsClient/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetSSMParameterValue(t *testing.T) {
	expectedParamValue := "secretToken123"
	parameterName := "testing-param1"
	var expectedErr error = nil

	ctrl := gomock.NewController(t)
	m := mocks.NewMockSSMAPI(ctrl)
	ctx := context.TODO()

	var fakeResponse ssm.GetParameterOutput = ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			Name:  &parameterName,
			Value: &expectedParamValue,
		},
	}

	m.EXPECT().GetParameter(gomock.Any(), gomock.Eq(
		&ssm.GetParameterInput{
			WithDecryption: true,
			Name:           &parameterName,
		})).Return(&fakeResponse, nil)

	paramValue, err := GetSSMParameterValue(ctx, m, parameterName)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedParamValue, paramValue, "Wrong paramValue. Expected %s, got %s", expectedParamValue, paramValue)
}

func TestGetSSMParameterValueError(t *testing.T) {
	expectedParamValue := ""
	parameterName := "testing-param1"
	var expectedErr error = fmt.Errorf("Can't find data")

	ctrl := gomock.NewController(t)
	m := mocks.NewMockSSMAPI(ctrl)
	ctx := context.TODO()

	m.EXPECT().GetParameter(gomock.Any(), gomock.Eq(
		&ssm.GetParameterInput{
			WithDecryption: true,
			Name:           &parameterName,
		})).Return(&ssm.GetParameterOutput{}, expectedErr)

	paramValue, err := GetSSMParameterValue(ctx, m, parameterName)
	assert.Equalf(t, err, expectedErr, "Wrong err. Expected %s, got %s", expectedErr, err)
	assert.Equalf(t, expectedParamValue, paramValue, "Wrong paramValue. Expected %s, got %s", expectedParamValue, paramValue)
}
