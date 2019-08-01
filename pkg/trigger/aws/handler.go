package aws

import (
	"time"

	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/trigger/aws/autoscaling"

	"github.com/golang/glog"
)

type HookHandler struct {
	AutoScaling *autoscaling.AutoScaling
	Drainer     drainer.Drainer
}

func (h *HookHandler) Loop(nodeName string) {
	glog.Infof("Running node drainer on node '%s' on instance '%s' in region '%s' and profile '%s'",
		nodeName, h.AutoScaling.Options.InstanceID, h.AutoScaling.Options.Region, h.AutoScaling.Options.Profile)
	for {
		glog.Infof("Sleeping %s seconds", h.AutoScaling.Options.LoopSleepTime)
		time.Sleep(h.AutoScaling.Options.LoopSleepTime)

		status, autoScalingGroupName, err := h.AutoScaling.GetInstanceStatusAndAutoScalingGroupName(&h.AutoScaling.Options.InstanceID)
		if err != nil {
			glog.Warningf("Can not get instance status and auto scaling group name, will try again: %s", err)
			continue
		}
		glog.Infof("Status of instance '%v' is '%v', autoscaling group is '%v'", h.AutoScaling.Options.InstanceID, *status, *autoScalingGroupName)
		if !h.AutoScaling.IsTerminating(status) {
			continue
		}

		err = h.Drainer.Drain(nodeName)
		if err != nil {
			glog.Warningf("Not all pods on this host can be evicted, will try again: %s", err)
			continue
		}
		glog.Infof("All evictable pods are gone, notifying AutoScalingGroup that instance '%v' can be shutdown", h.AutoScaling.Options.InstanceID)

		lifecycleHookName, err := h.AutoScaling.GetLifecycleHookName(autoScalingGroupName)
		if err != nil {
			glog.Warningf("Can not get lifecycle hook, will try again: %s", err)
			continue
		}

		glog.Infof("Sending notification to auto scaling group '%v' and lifecycle hook '%v'", *autoScalingGroupName, *lifecycleHookName)
		err = h.AutoScaling.SendNotification(&h.AutoScaling.Options.InstanceID, autoScalingGroupName, lifecycleHookName)
		if err != nil {
			glog.Warningf("Can not send notification, will try again: %s", err)
			continue
		}

		if h.AutoScaling.Options.ForceLoopBreak {
			glog.Warning("Reconciliation loop force-brake (normal only in tests)")
			break
		}
		glog.Infof("Sleeping %s, expecting that instance will be shut down in this time", h.AutoScaling.Options.ShutdownSleep)
		time.Sleep(h.AutoScaling.Options.ShutdownSleep)
	}
}
