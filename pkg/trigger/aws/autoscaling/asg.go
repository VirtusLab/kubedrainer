package autoscaling

import (
	"time"

	"github.com/VirtusLab/kubedrainer/pkg/drainer"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
)

const (
	// InstanceTerminatingStatus describes EC2 instance termination status
	InstanceTerminatingStatus = "Terminating:Wait"

	// LifecycleActionResultContinue describes ASG instance lifecycle continue result
	LifecycleActionResultContinue = "CONTINUE"
)

// AutoScaling type is a AWS EC2 AutoScaling API facade
type AutoScaling struct {
	AutoScaling autoscalingiface.AutoScalingAPI
	Options     *Options
}

type Options struct {
	InstanceID     string
	Region         string
	Profile        string
	LoopSleepTime  time.Duration
	ShutdownSleep  time.Duration
	ForceLoopBreak bool
}

type HookHandler struct {
	AutoScaling *AutoScaling
	Drainer     drainer.Drainer
}

func NewAutoScaling(session *session.Session, options *Options) *AutoScaling {
	return &AutoScaling{
		AutoScaling: autoscaling.New(session, aws.NewConfig().WithRegion(options.Region)),
		Options:     options,
	}
}

// GetInstanceStatusAndAutoScalingGroupName get an AWS EC2 instance status and its ASG name by instanceID
func (a *AutoScaling) GetInstanceStatusAndAutoScalingGroupName(instanceID *string) (*string, *string, error) {
	request := &autoscaling.DescribeAutoScalingInstancesInput{
		InstanceIds: []*string{instanceID},
	}

	result, err := a.AutoScaling.DescribeAutoScalingInstances(request)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}
	instanceCount := len(result.AutoScalingInstances)
	if instanceCount != 1 {
		return nil, nil, errors.Errorf("Expected exactly one instance, got: '%d'", instanceCount)
	}
	lifecycleState := result.AutoScalingInstances[0].LifecycleState
	autoScalingGroupName := result.AutoScalingInstances[0].AutoScalingGroupName
	return lifecycleState, autoScalingGroupName, nil
}

// GetLifecycleHookName gets an AWS ASG lifecycle hook name by autoScalingGroupName
func (a *AutoScaling) GetLifecycleHookName(autoScalingGroupName *string) (*string, error) {
	request := &autoscaling.DescribeLifecycleHooksInput{
		AutoScalingGroupName: autoScalingGroupName,
	}

	result, err := a.AutoScaling.DescribeLifecycleHooks(request)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	hooksCount := len(result.LifecycleHooks)
	if hooksCount != 1 {
		return nil, errors.Errorf("Expected exactly one lifecycle hook, got: '%d'", hooksCount)
	}
	return result.LifecycleHooks[0].LifecycleHookName, nil
}

// SendNotification sends a notification to AWS ASG using provided instanceID, autoScalingGroupName and lifecycleHookName
func (a *AutoScaling) SendNotification(instanceID *string, autoScalingGroupName *string, lifecycleHookName *string) error {
	input := &autoscaling.CompleteLifecycleActionInput{
		AutoScalingGroupName:  autoScalingGroupName,
		LifecycleActionResult: aws.String(LifecycleActionResultContinue),
		LifecycleHookName:     lifecycleHookName,
		InstanceId:            instanceID,
	}

	_, err := a.AutoScaling.CompleteLifecycleAction(input)
	return errors.Wrap(err)
}

// IsTerminating returns true if the provided status is in terminating state
func (a *AutoScaling) IsTerminating(status *string) bool {
	if status == nil {
		return false
	}
	return *status == InstanceTerminatingStatus
}
