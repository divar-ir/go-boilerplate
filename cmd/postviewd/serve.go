package main

import (
	"context"
	"fmt"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/adaptors"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/middlewares"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache/multilayercache"
	"github.com/allegro/bigcache"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/app/core"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/errors"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/sql"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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

	config := loadConfigOrPanic(cmd)

	configureLoggerOrPanic(config.Logging)

	prometheusMetricServer := startPrometheusMetricServerOrPanic(config.MetricListenPort)
	defer shutdownPrometheusMetricServerOrPanic(prometheusMetricServer)

	providerInstance := getProvider(config)
	cacheInstance := getCache(config)
	servicer := core.New(providerInstance, cacheInstance)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ListenPort))
	if err != nil {
		panicWithError(err, "failed to listen")
	}

	grpcServer := configureServer(config)
	postview.RegisterPostViewServer(grpcServer, servicer)

	serverCtx, serverCancel := makeServerCtx()
	defer serverCancel()
	var serverWaitGroup sync.WaitGroup

	serverWaitGroup.Add(1)
	go func() {
		defer serverWaitGroup.Done()

		if err := grpcServer.Serve(listener); err != nil {
			panicWithError(err, "failed to serve")
		}
	}()

	if err := declareReadiness(); err != nil {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	grpcServer.GracefulStop()

	serverWaitGroup.Wait()
}

func getProvider(config *Config) provider.PostProvider {
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

func getCache(config *Config) cache.Layer {
	var cacheLayers []cache.Layer
	if config.Cache.Redis.Enabled {
		redisClient := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", config.Cache.Redis.Host, config.Cache.Redis.Port),
			DB:   config.Cache.Redis.DB,
		})
		// Ping Redis
		err := redisClient.Ping().Err()
		if err != nil {
			panicWithError(err, "fail to connect to redis")
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
			panicWithError(err, "fail to initialize big cache")
		}

		cacheLayers = append(cacheLayers, adaptors.NewBigCacheAdaptor(bigCacheInstance))

	}

	cacheInstance := multilayercache.New(cacheLayers...)

	cacheInstance = middlewares.NewInstrumentationMiddleware(
		cacheInstance, cacheMetrics.With(map[string]string{
			"cache_type": "multilayer",
		}))

	return cacheInstance
}

func configureServer(config *Config) *grpc.Server {
	logEntry := logrus.WithFields(map[string]interface{}{
		"app": "postviewd",
	})

	interceptors := []grpc.UnaryServerInterceptor{
		grpc_logrus.UnaryServerInterceptor(logEntry),
		errors.UnaryServerInterceptor,
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(),
	}

	return grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...)))
}

func loadConfigOrPanic(cmd *cobra.Command) *Config {
	config, err := LoadConfig(cmd)
	if err != nil {
		panicWithError(err, "Failed to load configurations.")
	}
	return config
}

func configureLoggerOrPanic(loggerConfig LoggingConfig) {
	if err := configureLogging(&loggerConfig); err != nil {
		panicWithError(err, "Failed to configure logger.")
	}
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

func startPrometheusMetricServerOrPanic(listenPort int) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: promhttp.Handler(),
	}

	go listenAndServePrometheusMetrics(server)

	return server
}

func listenAndServePrometheusMetrics(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		panicWithError(err, "failed to start liveness http probe listener")
	}
}

func shutdownPrometheusMetricServerOrPanic(server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		panicWithError(err, "Failed to shutdown prometheus metric server")
	}
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
