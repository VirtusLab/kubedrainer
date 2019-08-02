package main

import (
	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/kubernetes"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func drainCmd(options *drainer.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "drain",
		Short: "Drain a node",
		Long:  `Drain a node by cordoning and pod eviction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeName := viper.GetString("node")
			glog.V(1).Infof("nodeName: '%s'", nodeName)

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
