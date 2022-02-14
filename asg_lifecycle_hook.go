package main

import (
	"encoding/json"
	"time"
)

type LifecycleMessage struct {
	Origin               string    `json:"Origin"`
	LifecycleHookName    string    `json:"LifecycleHookName"`
	Destination          string    `json:"Destination"`
	AccountId            string    `json:"AccountId"`
	RequestId            string    `json:"RequestId"`
	LifecycleTransition  string    `json:"LifecycleTransition"`
	AutoScalingGroupName string    `json:"AutoScalingGroupName"`
	Service              string    `json:"Service"`
	Time                 time.Time `json:"Time"`
	EC2InstanceId        string    `json:"EC2InstanceId"`
	NotificationMetadata map[string][]string
	LifecycleActionToken string `json:"LifecycleActionToken"`
}

func (obj *LifecycleMessage) UnmarshalJSON(data []byte) error {
	temp := struct {
		Origin                  string    `json:"Origin"`
		LifecycleHookName       string    `json:"LifecycleHookName"`
		Destination             string    `json:"Destination"`
		AccountId               string    `json:"AccountId"`
		RequestId               string    `json:"RequestId"`
		LifecycleTransition     string    `json:"LifecycleTransition"`
		AutoScalingGroupName    string    `json:"AutoScalingGroupName"`
		Service                 string    `json:"Service"`
		Time                    time.Time `json:"Time"`
		EC2InstanceId           string    `json:"EC2InstanceId"`
		NotificationMetadataRaw string    `json:"NotificationMetadata"`
		LifecycleActionToken    string    `json:"LifecycleActionToken"`
	}{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	obj.Origin = temp.Origin
	obj.LifecycleHookName = temp.LifecycleHookName
	obj.Destination = temp.Destination
	obj.AccountId = temp.AccountId
	obj.RequestId = temp.RequestId
	obj.LifecycleTransition = temp.LifecycleTransition
	obj.AutoScalingGroupName = temp.AutoScalingGroupName
	obj.Service = temp.Service
	obj.Time = temp.Time
	obj.EC2InstanceId = temp.EC2InstanceId
	obj.LifecycleActionToken = temp.LifecycleActionToken
	if err := json.Unmarshal([]byte(temp.NotificationMetadataRaw), &obj.NotificationMetadata); err != nil {
		return err
	}

	return nil
}
