package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/scheduler"
)

type Scheduler scheduler.CreateScheduleInput

type target struct {
	*scheduler.Target
	DeadLetterConfig            *scheduler.DeadLetterConfig            `json:"DeadLetterConfig,omitempty"`
	EcsParameters               *scheduler.EcsParameters               `json:"EcsParameters,omitempty"`
	EventBridgeParameters       *scheduler.EventBridgeParameters       `json:"EventBridgeParameters,omitempty"`
	Input                       *string                                `json:"Input,omitempty"`
	KinesisParameters           *scheduler.KinesisParameters           `json:"KinesisParameters,omitempty"`
	RetryPolicy                 *scheduler.RetryPolicy                 `json:"RetryPolicy,omitempty"`
	SageMakerPipelineParameters *scheduler.SageMakerPipelineParameters `json:"SageMakerPipelineParameters,omitempty"`
	SqsParameters               *scheduler.SqsParameters               `json:"SqsParameters,omitempty"`
}

type flexibleTimeWindow struct {
	*scheduler.FlexibleTimeWindow
	MaximumWindowInMinutes *int64 `json:"MaximumWindowInMinutes,omitempty"`
}

func (s *Scheduler) Register() string {
	target := target{
		Target: s.Target,
	}
	flexibleTimeWindow := flexibleTimeWindow{
		FlexibleTimeWindow: s.FlexibleTimeWindow,
	}
	targetJSON, _ := json.Marshal(target)
	flexibleTimeWindowJSON, _ := json.Marshal(flexibleTimeWindow)

	cmd := fmt.Sprintf("aws scheduler create-schedule --action-after-completion %s --schedule-expression \"at(%s)\" --name %s --target '%s' --flexible-time-window '%s'", *s.ActionAfterCompletion, *s.ScheduleExpression, *s.Name, targetJSON, flexibleTimeWindowJSON)
	return cmd
}

func (s *Scheduler) Deregister() string {
	cmd := fmt.Sprintf("aws scheduler delete-schedule --name %s", *s.Name)
	return cmd
}
