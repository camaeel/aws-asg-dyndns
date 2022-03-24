package awsClient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

type AutoscalingFakeClient struct {
	CompleteLifecycleActionOutput *autoscaling.CompleteLifecycleActionOutput
	CompleteLifecycleActionErr    error
}

func (c *AutoscalingFakeClient) CompleteLifecycleAction(ctx context.Context, params *autoscaling.CompleteLifecycleActionInput, optFns ...func(*autoscaling.Options)) (*autoscaling.CompleteLifecycleActionOutput, error) {
	return c.CompleteLifecycleActionOutput, c.CompleteLifecycleActionErr
}
