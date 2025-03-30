//nolint:forbidigo // it's okay to use fmt in this file
package config

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Users                       string `mapstructure:"users" default:""`
	JWTSecret                   string `mapstructure:"jwt_secret" default:""`
	Verbose                     bool   `mapstructure:"verbose" default:"false"`
	Port                        int    `mapstructure:"port" default:"8080"`
	DBPath                      string `mapstructure:"dbpath" default:":memory:"`
	Prefill                     bool   `mapstructure:"prefill" default:"false"`
	DisableImporters            bool   `mapstructure:"disableimporters" default:"false"`
	DisableCurrenciesRatesFetch bool   `mapstructure:"disablecurrenciesratesfetch" default:"false"`
}

func InitiateConfig(cfgFile string) (*Config, error) {
	cfg := Config{}

	setDefaultsFromStruct(&cfg)
	viper.SetEnvPrefix("GB")
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.Verbose {
		fmt.Printf("Config: %+v\n", cfg)
	}

	return &cfg, nil
}

func setDefaultsFromStruct(s interface{}) {
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		defaultValue := field.Tag.Get("default")
		viper.SetDefault(field.Name, defaultValue)
	}
}
