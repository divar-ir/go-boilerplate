package main

import (
	"strings"
	"time"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/errors"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/sql"
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
	Cache            CacheConfig
}

type CacheConfig struct {
	Redis    RedisConfig
	BigCache BigCacheConfig
}

type RedisConfig struct {
	Enabled        bool
	Host           string
	Port           int
	DB             int
	Prefix         string
	ExpirationTime time.Duration
}

type BigCacheConfig struct {
	Enabled            bool
	ExpirationTime     time.Duration
	MaxSpace           int
	Shards             int
	LifeWindow         time.Duration
	MaxEntriesInWindow int
	MaxEntrySize       int
	Verbose            bool
	HardMaxCacheSize   int
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
	viper.SetDefault("cache.redis.host", "127.0.0.1")
	viper.SetDefault("cache.redis.port", 6379)
	viper.SetDefault("cache.redis.db", 0)
	viper.SetDefault("cache.redis.expirationTime", 3*time.Hour)
	viper.SetDefault("cache.redis.prefix", "POST_VIEW")
	viper.SetDefault("cache.bigCache.shards", 1024)
	viper.SetDefault("cache.bigCache.maxEntriesInWindow", 1100*10*60)
	viper.SetDefault("cache.bigCache.maxEntrySize", 500)
	viper.SetDefault("cache.bigCache.verbose", true)
	viper.SetDefault("cache.bigCache.hardMaxCacheSize", 125)

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

func provideConfig(cmd *cobra.Command) (*Config, error) {
	config, err := LoadConfig(cmd)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load configurations.")
	}
	return config, nil
}
