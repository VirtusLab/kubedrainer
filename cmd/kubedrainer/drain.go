package main

import (
	"time"

	"github.com/VirtusLab/kubedrainer/internal/settings"
	"github.com/VirtusLab/kubedrainer/internal/stringer"
	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/kubernetes"
	"github.com/rs/zerolog/log"

	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// DrainOptions holds the drainer options
type DrainOptions struct {
	Kubernetes *kubernetes.Options
	Drainer    *drainer.Options
}

// DrainFlags holds the drainer flags
type DrainFlags struct {
	Kubernetes *pflag.FlagSet
	Drainer    *pflag.FlagSet
}

func drainCmd() *cobra.Command {
	options := &DrainOptions{
		Kubernetes: kubernetes.DefaultOptions(),
		Drainer: &drainer.Options{
			GracePeriodSeconds:  -1,
			Timeout:             60 * time.Second,
			DeleteLocalData:     true,
			IgnoreAllDaemonSets: true,
		},
	}

	flags := &DrainFlags{
		Kubernetes: kubernetesFlags(options.Kubernetes),
		Drainer:    drainerFlags(options.Drainer),
	}

	cmd := &cobra.Command{
		Use:   "drain",
		Short: "Drain a node",
		Long:  `Drain a node by cordoning and pod eviction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			setupLogging()
			log.Info().Msg("Running locally")

			if err := options.Parse(cmd); err != nil {
				return err
			}

			client, err := kubernetes.New(options.Kubernetes)
			if err != nil {
				return err
			}

			if len(options.Drainer.Node) == 0 {
				return errors.New("No node name provided")
			}

			d := drainer.New(client, options.Drainer)
			err = d.Drain(options.Drainer.Node)
			return errors.Wrap(err)
		},
	}

	flags.AddTo(cmd.PersistentFlags())
	return cmd
}

// AddTo adds the flags to the given flag set
func (f *DrainFlags) AddTo(flags *pflag.FlagSet) {
	flags.AddFlagSet(f.Kubernetes)
	flags.AddFlagSet(f.Drainer)
}

// Parse parses all flags and settings to options
func (o *DrainOptions) Parse(cmd *cobra.Command) error {
	settings.Bind(cmd.Flags()) // needs to be run inside the command and before any viper usage for flags to be visible

	if debug {
		log.Debug().Msgf("All keys: %+v", viper.AllKeys())
		log.Debug().Msgf("All settings: %+v", viper.AllSettings())
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			log.Debug().Msgf("'%s' -> flag: '%+v' | setting: '%+v'", flag.Name, flag.Value, viper.Get(flag.Name))
		})
		log.Debug().Msgf("Settings: %+v", *o)
	}

	if err := settings.Parse(o.Kubernetes); err != nil {
		return err
	}
	if err := settings.Parse(o.Drainer); err != nil {
		return err
	}
	log.Debug().Msgf("Parsed settings: %+v", *o)
	return nil
}

// String implements Stringer
func (o *DrainOptions) String() string {
	return stringer.Stringify(o)
}
