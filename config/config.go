package config

import (
	"github.com/anzx/pkg/jsontime"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
	Age int
	Banana bool
}

func flags() *flag.FlagSet {
	f := flag.NewFlagSet("Test", flag.PanicOnError)

	f.StringP("name", "n", "DefaultName", "")
	f.IntP("age", "a", 0, "")
	f.BoolP("banana", "b", false, "")

	return f
}

// priority: viper -> env -> flags
func configure(v *viper.Viper, f *flag.FlagSet, e [][]string) (*Config, error) {
	v.SetConfigFile("example.yaml")

	// Process config file
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Process environment variables array
	v.AutomaticEnv()

	// Process flags
	err = v.BindPFlags(f)
	if err != nil {
		panic(err)
	}

	var config *Config

	err = v.Unmarshal(&config, viper.DecodeHook(jsontime.DurationMapstructureDecodeHookFunc))
	if err != nil {
		panic(err)
	}

	return config, nil
}