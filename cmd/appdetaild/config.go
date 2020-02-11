package main

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Config the application's configuration structure
type Config struct {
	Logging          LoggingConfig
	ConfigFile       string
	ListenPort       int
}


// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	// Setting defaults for this application
	viper.SetDefault("logging.SentryEnabled", false)
	viper.SetDefault("logging.level", "error")
	viper.SetDefault("listenPort", 8080)

	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("APPDETAIL")
	viper.AutomaticEnv()

	// Read Config from Flags
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	// Read Config from file
	if configFile, err := cmd.Flags().GetString("config-file"); err == nil && configFile != "" {
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func provideConfig(cmd *cobra.Command) (*Config, error) {
	config, err := LoadConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load configurations.")
	}
	return config, nil
}
