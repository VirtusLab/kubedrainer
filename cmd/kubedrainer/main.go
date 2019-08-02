package main

import (
	"flag"
	"reflect"
	"time"

	"github.com/VirtusLab/kubedrainer/internal/version"
	"github.com/VirtusLab/kubedrainer/pkg/drainer"
	"github.com/VirtusLab/kubedrainer/pkg/trigger/aws/autoscaling"

	"github.com/golang/glog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var appName = "kubedrainer"
var cfgFile string

var kubeConfigFlags *genericclioptions.ConfigFlags
var drainerOptions = &drainer.Options{
	GracePeriodSeconds:  -1,
	Timeout:             60 * time.Second,
	DeleteLocalData:     true,
	IgnoreAllDaemonSets: true,
}
var awsOptions = &autoscaling.Options{
	LoopSleepTime: 10 * time.Second,
	ShutdownSleep: 6 * time.Minute,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "Kubernetes Node Drainer",
	Long:    `Kubernetes Node Drainer helps to evicts pods form node before shutdown`,
	Version: version.Long(),
}

func init() {
	// initialization actions
	cobra.OnInitialize(
		initConfig,
	)

	// add global flags
	addGlogFlags(pflag.CommandLine)

	rootFlags := rootCmd.PersistentFlags()
	rootFlags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubedrainer.yaml)")

	// add drain command
	drainCmd := drainCmd(drainerOptions)
	drainFlags := drainCmd.PersistentFlags()
	addKubeConfigFlags(drainFlags)
	addDrainerFlags(drainFlags, drainerOptions)
	rootCmd.AddCommand(drainCmd)

	// add serve command
	serveCmd := serveCmd(drainerOptions, awsOptions)
	serveFlags := serveCmd.PersistentFlags()
	addKubeConfigFlags(serveFlags)
	addAutoscalingFlags(serveFlags, awsOptions)
	addDrainerFlags(serveFlags, drainerOptions)
	rootCmd.AddCommand(serveCmd)
}

func main() {
	// make sure we always get logs
	defer glog.Flush()

	// Adds all child commands to the root command and sets flags appropriately.
	// This is called by main.main(). It only needs to happen once to the rootCmd.
	if err := rootCmd.Execute(); err != nil {
		glog.V(1).Infof("Stack trace (%s): %+v", reflect.TypeOf(err), err)
		glog.Exitln(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	glog.V(3).Info("initConfig")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			glog.Errorln(err)
		}

		// Search config in home directory with name ".kubedrainer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + appName)
	}

	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv() // read environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.V(1).Info("No config file found.")
		} else {
			glog.Errorf("Config file found, but cannot be read.")
		}
	} else {
		glog.Infof("Using config file: '%s'", viper.ConfigFileUsed())
	}
}

func addGlogFlags(flags *pflag.FlagSet) {
	// the following line exists to make glog happy, for more information
	// see: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	_ = flag.CommandLine.Parse([]string{})

	// set glog defaults
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("alsologtostderr", "false")

	// add glog flags to cobra
	flags.AddGoFlagSet(flag.CommandLine)
}

func addDrainerFlags(flags *pflag.FlagSet, options *drainer.Options) {
	flags.BoolVar(&options.DryRun, "dry-run", options.DryRun, "If true, only print the object that would be sent, without sending it.")
	flags.BoolVar(&options.Force, "force", options.Force, "Continue even if there are pods not managed by a ReplicationController, ReplicaSet, Job, DaemonSet or StatefulSet.")
	flags.BoolVar(&options.IgnoreAllDaemonSets, "ignore-daemonsets", options.IgnoreAllDaemonSets, "Ignore DaemonSet-managed pods.")
	flags.BoolVar(&options.DeleteLocalData, "delete-local-data", options.DeleteLocalData, "Continue even if there are pods using emptyDir (local data that will be deleted when the node is drained).")
	flags.IntVar(&options.GracePeriodSeconds, "grace-period", options.GracePeriodSeconds, "Period of time in seconds given to each pod to terminate gracefully. If negative, the default value specified in the pod will be used.")
	flags.DurationVar(&options.Timeout, "timeout", options.Timeout, "The length of time to wait before giving up, zero means infinite")
	flags.StringVarP(&options.Selector, "selector", "l", options.Selector, "Selector (label query) to filter on")
	flags.StringVarP(&options.PodSelector, "pod-selector", "", options.PodSelector, "Label selector to filter pods on the node")

	flags.String("node", "", "Kubernetes node name to drain")
	_ = viper.BindPFlag("node", flags.Lookup("node"))
	_ = viper.BindEnv("node")
}

func addKubeConfigFlags(flags *pflag.FlagSet) {
	// integrate with kubeconfig
	kubeConfigFlags = genericclioptions.NewConfigFlags(true)
	kubeConfigFlags.Namespace = nil        // disable 'namespace' flag
	kubeConfigFlags.Impersonate = nil      // disable 'as' flag
	kubeConfigFlags.ImpersonateGroup = nil // disable 'as-group' flag
	kubeConfigFlags.AddFlags(flags)        // add all kubeconfig specific flags
}

func addAutoscalingFlags(flags *pflag.FlagSet, options *autoscaling.Options) {
	flags.StringVar(&options.Region, "region", options.Region, "AWS Region to use")
	flags.StringVar(&options.Profile, "profile", options.Region, "AWS Profile to use")
	flags.StringVar(&options.InstanceID, "instance-id", options.InstanceID, "AWS EC2 instance ID to terminate")
}
