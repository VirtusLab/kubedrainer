package aws

import (
	"time"

	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/trigger/aws/autoscaling"
	"github.com/rs/zerolog/log"
)

// HookHandler implements a drainer trigger for AWS ASG instance lifecycle hook
type HookHandler struct {
	AutoScaling *autoscaling.AutoScaling
	Drainer     drainer.Drainer
}

// Loop starts an infinite handler loop
func (h *HookHandler) Loop(nodeName string) {
	var drained bool
	log.Info().Msgf("Running node drainer on node '%s' on instance '%s' in region '%s' and profile '%s'",
		nodeName, h.AutoScaling.Options.InstanceID, h.AutoScaling.Options.Region, h.AutoScaling.Options.Profile)
	for {
		log.Info().Msgf("Sleeping %s seconds", h.AutoScaling.Options.LoopSleepTime)
		time.Sleep(h.AutoScaling.Options.LoopSleepTime)

		status, autoScalingGroupName, err := h.AutoScaling.GetInstanceStatusAndAutoScalingGroupName(&h.AutoScaling.Options.InstanceID)
		if err != nil {
			log.Warn().Msgf("Can not get instance status and auto scaling group name, will try again: %s", err)
			continue
		}
		log.Info().Msgf("Status of instance '%v' is '%v', autoscaling group is '%v'", h.AutoScaling.Options.InstanceID, *status, *autoScalingGroupName)
		if !h.AutoScaling.IsTerminating(status) && !h.AutoScaling.IsTerminatingWait(status) {
			continue
		}

		if !drained {
			err = h.Drainer.Drain(nodeName)
			if err != nil {
				log.Warn().Msgf("Not all pods on this host can be evicted, will try again: %s", err)
				continue
			}
			drained = true
			log.Warn().Msg("All evictable pods are gone, waiting to enter Terminating:Wait state")
		}

		if !h.AutoScaling.IsTerminatingWait(status) {
			continue
		}

		log.Info().Msgf("Notifying AutoScalingGroup that instance '%v' can be shutdown", h.AutoScaling.Options.InstanceID)
		lifecycleHookName, err := h.AutoScaling.GetLifecycleHookName(autoScalingGroupName)
		if err != nil {
			log.Warn().Msgf("Can not get lifecycle hook, will try again: %s", err)
			continue
		}

		log.Info().Msgf("Sending notification to auto scaling group '%v' and lifecycle hook '%v'", *autoScalingGroupName, *lifecycleHookName)
		err = h.AutoScaling.SendNotification(&h.AutoScaling.Options.InstanceID, autoScalingGroupName, lifecycleHookName)
		if err != nil {
			log.Warn().Msgf("Can not send notification, will try again: %s", err)
			continue
		}

		if h.AutoScaling.Options.ForceLoopBreak {
			log.Warn().Msg("Reconciliation loop force-brake (normal only in tests)")
			break
		}
		log.Info().Msgf("Sleeping %s, expecting that instance will be shut down in this time", h.AutoScaling.Options.ShutdownSleep)
		time.Sleep(h.AutoScaling.Options.ShutdownSleep)
	}
}
