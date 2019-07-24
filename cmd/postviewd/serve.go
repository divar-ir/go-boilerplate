package main

import (
	"context"
	"fmt"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/grpcserver"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/metrics/prometheus"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/adaptors"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/middlewares"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/multilayercache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/errors"
	"github.com/allegro/bigcache"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/sql"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start Server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	printVersion()

	serverCtx, serverCancel := makeServerCtx()
	defer serverCancel()

	server, err := CreateServer(serverCtx, cmd)
	panicWithError(err, "failed to create server")

	var serverWaitGroup sync.WaitGroup

	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()

		if err := server.Serve(); err != nil {
			panicWithError(err, "failed to serve")
		}
	}()

	if err := declareReadiness(); err != nil {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	server.Stop()

	serverWaitGroup.Wait()
}

func provideServer(server postview.PostViewServer, config *Config, logger *logrus.Logger) (*grpcserver.Server, error) {
	return grpcserver.New(server, logger, config.ListenPort)
}

func provideProvider(config *Config, logger *logrus.Logger, prometheusMetric *prometheus.Server) provider.PostProvider {
	db, err := sql.GetDatabase(config.Database)
	if err != nil {
		logrus.WithError(err).WithField(
			"database", config.Database).Panic("failed to connect to DB")
		return nil
	}

	providerInstance := provider.NewSQL(db)
	providerInstance = provider.NewInstrumentationMiddleware(
		providerInstance, postProviderMetrics.With(map[string]string{
			"provider_type": "postgres",
		}))

	return providerInstance
}

func provideCache(config *Config, prometheusMetric *prometheus.Server) (cache.Layer, error) {
	var cacheLayers []cache.Layer
	if config.Cache.Redis.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", config.Cache.Redis.Host, config.Cache.Redis.Port),
			DB:   config.Cache.Redis.DB,
		})
		// Ping Redis
		err := redisClient.Ping().Err()
		if err != nil {
			return nil, errors.Wrap(err, "fail to connect to redis")
		}
		cacheLayers = append(cacheLayers, adaptors.NewRedisAdaptor(config.Cache.Redis.ExpirationTime, redisClient))

	}

	if config.Cache.BigCache.Enabled {
		bigCacheInstance, err := bigcache.NewBigCache(bigcache.Config{
			Shards:             config.Cache.BigCache.Shards,
			LifeWindow:         config.Cache.BigCache.ExpirationTime,
			MaxEntriesInWindow: config.Cache.BigCache.MaxEntriesInWindow,
			MaxEntrySize:       config.Cache.BigCache.MaxEntrySize,
			Verbose:            config.Cache.BigCache.Verbose,
			HardMaxCacheSize:   config.Cache.BigCache.HardMaxCacheSize,
		})
		if err != nil {
			return nil, errors.Wrap(err, "fail to initialize big cache")
		}

		cacheLayers = append(cacheLayers, adaptors.NewBigCacheAdaptor(bigCacheInstance))

	}

	cacheInstance := multilayercache.New(cacheLayers...)

	cacheInstance = middlewares.NewInstrumentationMiddleware(
		cacheInstance, cacheMetrics.With(map[string]string{
			"cache_type": "multilayer",
		}))

	return cacheInstance, nil
}

func makeServerCtx() (context.Context, context.CancelFunc) {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-gracefulStop
		cancel()
	}()

	return ctx, cancel
}

func declareReadiness() error {
	// nolint: gosec
	file, err := os.Create("/tmp/readiness")
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer file.Close()

	_, err = file.WriteString("ready")
	return err
}

func panicWithError(err error, format string, args ...interface{}) {
	logrus.WithError(err).Panicf(format, args...)
}
