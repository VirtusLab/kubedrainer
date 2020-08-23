package kubernetes

import (
	"github.com/VirtusLab/kubedrainer/internal/stringer"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/rs/zerolog/log"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

// Options is an alias for kubectl "ConfigFlags"
type Options struct {
	*genericclioptions.ConfigFlags
}

// Client is an alias for kubernetes "Clientset"
type Client = kubernetes.Clientset

// DefaultOptions creates a default Options instance
func DefaultOptions() *Options {
	return &Options{
		genericclioptions.NewConfigFlags(true),
	}
}

// New returns a Kubernetes API client using kubeconfig
func New(options *Options) (*Client, error) {
	clientConfig, err := options.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	log.Debug().Msgf("Context: %s", *options.Context)
	log.Debug().Msgf("Configured Host: %s", clientConfig.Host)
	log.Debug().Msgf("Configured AuthProvider: %s", clientConfig.AuthProvider)
	log.Debug().Msgf("Configured ExecProvider: %s", clientConfig.ExecProvider)

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	version, err := client.ServerVersion()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	log.Debug().Msgf("Server version: %s", version.String())

	return client, err
}

// String implements Stringer
func (o *Options) String() string {
	return stringer.Stringify(o)
}
