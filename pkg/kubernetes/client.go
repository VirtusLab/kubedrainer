package kubernetes

import (
	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/golang/glog"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// Options is an alias for kubectl "ConfigFlags"
type Options = genericclioptions.ConfigFlags

// Client is an alias for kubernetes "Clientset"
type Client = kubernetes.Clientset

// New returns a Kubernetes API client using kubeconfig
func New(options *Options) (*Client, error) {
	clientConfig, err := options.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if glog.V(2) {
		glog.Infof("Context: %s", *options.Context)
	}
	if glog.V(4) {
		glog.Infof("Configured Host: %s", clientConfig.Host)
		glog.Infof("Configured AuthProvider: %s", clientConfig.AuthProvider)
		glog.Infof("Configured ExecProvider: %s", clientConfig.ExecProvider)
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if glog.V(1) {
		version, err := client.ServerVersion()
		if err != nil {
			return nil, errors.Wrap(err)
		}
		glog.Infof("Server version: %s", version.String())
	}

	return client, err
}
