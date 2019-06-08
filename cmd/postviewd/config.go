package main

import (
	"strings"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Config the application's configuration structure
type Config struct {
	Logging          LoggingConfig
	ConfigFile       string
	ListenPort       int
	MetricListenPort int
	Database         sql.PostgresConfig
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	// Setting defaults for this application
	viper.SetDefault("logging.SentryEnabled", false)
	viper.SetDefault("logging.level", "error")
	viper.SetDefault("listenPort", 8080)
	viper.SetDefault("metricListenPort", 8081)
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "postview")
	viper.SetDefault("database.ssl", false)
	viper.SetDefault("database.maxIdleConnection", 0)
	viper.SetDefault("database.maxOpenConnection", 0)

	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("POSTVIEW")
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
