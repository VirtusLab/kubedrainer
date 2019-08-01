package kubernetes

import (
	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/golang/glog"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// Client returns a Kubernetes API client using kubeconfig
func Client(kubeConfigFlags *genericclioptions.ConfigFlags) (*kubernetes.Clientset, error) {
	clientConfig, err := kubeConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if glog.V(3) {
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
