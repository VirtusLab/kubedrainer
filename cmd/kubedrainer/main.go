package main

import (
	"reflect"

	"github.com/VirtusLab/kubedrainer/internal/version"
	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/kubernetes"
	"github.com/VirtusLab/kubedrainer/pkg/trigger/aws/autoscaling"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var appName = "kubedrainer"
var cfgFile string
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "Kubernetes Node Drainer",
	Long:    `Kubernetes Node Drainer helps to evicts pods form node before shutdown`,
	Version: version.Long(),
}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cobra.OnInitialize(
		initConfig,
	)

	rootFlags := rootCmd.PersistentFlags()
	rootFlags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubedrainer.yaml)")
	rootFlags.BoolVarP(&debug, "debug", "d", false, "sets log level to debug")

	rootCmd.AddCommand(drainCmd())
	rootCmd.AddCommand(serveCmd())
}

func main() {
	// handle config auto-reload
	viper.OnConfigChange(func(e fsnotify.Event) {
		// TODO
		log.Warn().Msg("Config auto reload not implemented!")
	})

	// Adds all child commands to the root command and sets flags appropriately.
	// This is called by main.main(). It only needs to happen once to the rootCmd.
	if err := rootCmd.Execute(); err != nil {
		exit(err)
	}
}

func exit(err error) {
	log.Debug().Msgf("Stack trace (%s): %+v", reflect.TypeOf(err), err)
	log.Fatal().Msgf("%v", err)
}

func setupLogging() {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error().Msgf("%v", err)
		}

		// Search config in home directory with name ".kubedrainer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + appName)
	}

	viper.SetEnvPrefix(appName)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn().Msg("No config file found")
		} else {
			log.Error().Msg("Config file found, but cannot be read")
		}
	} else {
		log.Info().Msgf("Using config file: '%s'", viper.ConfigFileUsed())
		viper.WatchConfig()
	}
}

func drainerFlags(options *drainer.Options) *pflag.FlagSet {
	var flags = pflag.NewFlagSet("drainer", pflag.ContinueOnError)
	flags.String("node", options.Node, "Kubernetes node name to drain")
	flags.Bool("force", options.Force, "Continue even if there are pods not managed by a ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet.")
	flags.Bool("ignore-daemonsets", options.IgnoreAllDaemonSets, "Ignore DaemonSet-managed pods.")
	flags.Bool("delete-local-data", options.DeleteLocalData, "Continue even if there are pods using emptyDir (local data that will be deleted when the node is drained).")
	flags.Int("grace-period", options.GracePeriodSeconds, "Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used.")
	flags.Duration("timeout", options.Timeout, "The length of time to wait before giving up, zero means infinite")
	flags.StringP("selector", "l", options.Selector, "Selector (label query) to filter on")
	flags.StringP("pod-selector", "", options.PodSelector, "Label selector to filter pods on the node")
	return flags
}

func kubernetesFlags(options *kubernetes.Options) *pflag.FlagSet {
	var flags = pflag.NewFlagSet("kubernetes", pflag.ContinueOnError)
	// integrate with kubeconfig
	options.Namespace = nil        // disable 'namespace' flag
	options.Impersonate = nil      // disable 'as' flag
	options.ImpersonateGroup = nil // disable 'as-group' flag
	options.AddFlags(flags)        // add all kubeconfig specific flags
	return flags
}

func autoscalingFlags(options *autoscaling.Options) *pflag.FlagSet {
	var flags = pflag.NewFlagSet("autoscaling", pflag.ContinueOnError)
	flags.String("region", options.Region, "AWS Region to use")
	flags.String("profile", options.Region, "AWS Profile to use")
	flags.String("instance-id", options.InstanceID, "AWS EC2 instance ID to terminate")
	return flags
}
