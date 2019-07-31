package main

import (
	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/VirtusLab/go-extended/pkg/matcher"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/kubernetes"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/kubernetes/node"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/trigger/aws"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/trigger/aws/autoscaling"
	"github.com/VirtusLab/kubedrainer/cmd/pkg/trigger/aws/metadata"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	k8s "k8s.io/client-go/kubernetes"
)

// serveCmd represents the serve command
func serveCmd(drainerOptions *drainer.Options, asgOptions *autoscaling.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run node drainer as server",
		Long:  `Run node drainer as server with the provided configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("")
			}

			glog.Info("Running as server")

			nodeName := args[0]
			glog.V(1).Infof("NodeName (args[0]): '%s'", nodeName)
			glog.V(1).Infof("Instance ID: '%s'", asgOptions.InstanceID)
			glog.V(1).Infof("Region: '%s'", asgOptions.Region)

			kubernetesClient, err := kubernetes.Client(kubeConfigFlags)
			if err != nil {
				return err
			}

			awsSession, _, err := aws.SessionConfig(asgOptions.Region, asgOptions.Profile)
			if err != nil {
				return err
			}

			// get information from Kubernetes API if necessary
			if len(asgOptions.Region) == 0 && len(asgOptions.InstanceID) == 0 {
				glog.V(1).Info("Getting node information")
				region, instanceID, err := GetNodeInformation(nodeName, kubernetesClient)
				if err != nil {
					return err
				}
				if len(asgOptions.InstanceID) == 0 {
					asgOptions.InstanceID = instanceID
				} else {
					glog.V(1).Infof("Ignoring instance ID from node info '%s', using current: '%s'",
						instanceID, asgOptions.InstanceID)
				}
				if len(asgOptions.Region) == 0 {
					asgOptions.Region = region
				} else {
					glog.V(1).Infof("Ignoring region from node info '%s', using current: '%s'",
						instanceID, asgOptions.Region)
				}
			}

			// get information from AWS API if necessary
			if len(asgOptions.Region) == 0 && len(asgOptions.InstanceID) == 0 {
				glog.V(1).Info("Getting EC2 metadata")
				region, instanceID, err := GetMetadata(awsSession)
				if err != nil {
					return err
				}
				if len(asgOptions.InstanceID) == 0 {
					asgOptions.InstanceID = instanceID
				} else {
					glog.V(1).Infof("Ignoring instance ID from metadata '%s', using current: '%s'",
						instanceID, asgOptions.InstanceID)
				}
				if len(asgOptions.Region) == 0 {
					asgOptions.Region = region
				} else {
					glog.V(1).Infof("Ignoring region from metadata '%s', using current: '%s'",
						instanceID, asgOptions.Region)
				}
			}

			if len(asgOptions.Profile) == 0 {
				glog.V(1).Infof("Using default AWS API credentials profile")
				asgOptions.Profile = aws.DefaultProfile
			}

			// check preconditions
			if len(asgOptions.InstanceID) == 0 {
				return errors.New("No instance ID from flags, configuration, or metadata")
			}
			if len(asgOptions.Region) == 0 {
				return errors.New("No region from flags, configuration, or metadata")
			}
			if len(asgOptions.Profile) == 0 {
				return errors.New("No profile provided")
			}

			t := aws.HookHandler{
				Drainer:     drainer.New(kubernetesClient, drainerOptions),
				AutoScaling: autoscaling.NewAutoScaling(awsSession, asgOptions),
			}

			t.Loop(nodeName)

			return errors.Wrap(err)
		},
	}
}

func GetNodeInformation(nodeName string, kubernetesClient k8s.Interface) (string, string, error) {
	var region string
	var instanceID string

	n := &node.Node{
		Client: kubernetesClient,
	}
	providerName, providerSpecificNodeID, err := n.GetProviderId(nodeName)
	if err != nil {
		return "", "", err
	}
	switch providerName {
	case "aws":
		awsNodeIDExpression := `^/(?P<Region>[a-zA-Z0-9-]+)[a-z]/(?P<InstanceID>[a-zA-Z0-9-]+)$`
		results, ok := matcher.Must(awsNodeIDExpression).MatchGroups(providerSpecificNodeID)
		if !ok {
			return "", "", errors.Errorf("Can't match expression '%s' to '%s'",
				awsNodeIDExpression, providerSpecificNodeID)
		}
		region, ok = results["Region"]
		if !ok {
			return "", "", errors.Errorf("Missing 'Region' when expression '%s' was applied to '%s'",
				awsNodeIDExpression, providerSpecificNodeID)
		}
		instanceID, ok = results["InstanceID"]
		if !ok {
			return "", "", errors.Errorf("Missing 'InstanceID' when expression '%s' was applied to '%s'",
				awsNodeIDExpression, providerSpecificNodeID)
		}
	default:
		return "", "", errors.Errorf("Cloud provider not supported: '%s'", providerName)
	}

	return region, instanceID, nil
}

func GetMetadata(awsSession *session.Session) (string, string, error) {
	var region string
	var instanceID string
	m := metadata.New(awsSession)
	instanceID, region, err := m.GetCurrentInstanceIDAndRegion()
	switch err := err.(type) {
	case nil: // nothing
	case awserr.Error:
		if err.Code() == "EC2MetadataRequestError" {
			glog.Warning("No EC2 metadata available")
		} else {
			return "", "", errors.Wrap(err)
		}
	default:
		return "", "", errors.Wrap(err)
	}
	return region, instanceID, nil
}
