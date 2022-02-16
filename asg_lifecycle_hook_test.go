package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLifecycleMessageUnmarshal(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	input := "{\"Origin\":\"EC2\",\"LifecycleHookName\":\"aws-asg-dynds-create\",\"Destination\":\"AutoScalingGroup\",\"AccountId\":\"123123\",\"RequestId\":\"123123-123-abcd-aaa1-123123ASDZZZ\",\"LifecycleTransition\":\"autoscaling:EC2_INSTANCE_LAUNCHING\",\"AutoScalingGroupName\":\"test\",\"Service\":\"AWS Auto Scaling\",\"Time\":\"" + (now.UTC().Format("2006-01-02T15:04:05Z07:00")) + "\",\"EC2InstanceId\":\"i-123abcd0123\",\"NotificationMetadata\":\"{\\\"domainList\\\":[\\\"test.example.com\\\"]}\",\"LifecycleActionToken\":\"123123a-1234-abcd-def1-123123123\"}"

	var expected LifecycleMessage = LifecycleMessage{
		Origin:               "EC2",
		LifecycleHookName:    "aws-asg-dynds-create",
		Destination:          "AutoScalingGroup",
		AccountId:            "123123",
		RequestId:            "123123-123-abcd-aaa1-123123ASDZZZ",
		LifecycleTransition:  "autoscaling:EC2_INSTANCE_LAUNCHING",
		AutoScalingGroupName: "test",
		Service:              "AWS Auto Scaling",
		Time:                 now.UTC(),
		EC2InstanceId:        "i-123abcd0123",
		NotificationMetadata: map[string][]string{"domainList": {"test.example.com"}},
		LifecycleActionToken: "123123a-1234-abcd-def1-123123123"}

	var result LifecycleMessage

	err := json.Unmarshal([]byte(input), &result)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, expected, result, "Result should have the same value as expected")

}

func TestLifecycleMessageUnmarshalTest(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	input := "{\"AccountId\":\"123123123\",\"RequestId\":\"123123-asda-wwww-bfce-123123123\",\"AutoScalingGroupARN\":\"arn:aws:autoscaling:eu-north-1:123123:autoScalingGroup:123123-123123-123-123-123123123:autoScalingGroupName/testing\",\"AutoScalingGroupName\":\"testing\",\"Service\":\"AWS Auto Scaling\",\"Event\":\"autoscaling:TEST_NOTIFICATION\",\"Time\":\"" + (now.UTC().Format("2006-01-02T15:04:05Z07:00")) + "\"}"

	var expected LifecycleMessage = LifecycleMessage{
		// Destination: "AutoScalingGroup",
		AccountId: "123123123",
		RequestId: "123123-asda-wwww-bfce-123123123",
		// Event:  							"autoscaling:TEST_NOTIFICATION",
		AutoScalingGroupName: "testing",
		// AutoScalingGroupARN:  "arn:aws:autoscaling:eu-north-1:123123:autoScalingGroup:123123-123123-123-123-123123123:autoScalingGroupName/testing",
		Service: "AWS Auto Scaling",
		Time:    now.UTC(),
	}

	var result LifecycleMessage

	err := json.Unmarshal([]byte(input), &result)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, expected, result, "Result should have the same value as expected")

}
