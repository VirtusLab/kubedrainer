package settings

import (
	"github.com/VirtusLab/go-extended/pkg/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Bind binds given flags to settings
func Bind(flags *pflag.FlagSet) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(flag.Name, flag); err != nil {
			panic(err)
		}
		if err := viper.BindEnv(flag.Name); err != nil {
			panic(err)
		}
	})
}

// Parse parses (unmarshals) settings to a given interface
func Parse(target interface{}) error {
	err := viper.Unmarshal(target, func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = false
		config.ZeroFields = false
	})
	return errors.Wrap(err)
}
