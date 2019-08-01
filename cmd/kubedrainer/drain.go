package main

import (
	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/kubernetes"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func drainCmd(options *drainer.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "drain",
		Short: "Drain a node",
		Long:  `Drain a node by cordoning and pod eviction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("")
			}

			nodeName := args[0]
			glog.V(3).Infof("nodeName(args[0])=%v", nodeName)

			client, err := kubernetes.Client(kubeConfigFlags)
			if err != nil {
				return err
			}
			d := drainer.New(client, options)
			err = d.Drain(nodeName)
			return errors.Wrap(err)
		},
	}
}
