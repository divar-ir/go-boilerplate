package main

import (
	"fmt"
	"log"
	"net"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"
	"github.com/spf13/cobra"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/app/core"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"google.golang.org/grpc"
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

	providerInstance := provider.NewMemory()
	cacheInstance := cache.NewMemory()
	servicer := core.New(providerInstance, cacheInstance)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ListenPort))
	if err != nil {
		panicWithError(err, "failed to listen")
	}

	grpcServer := grpc.NewServer()
	postview.RegisterPostViewServer(grpcServer, servicer)

	if err := grpcServer.Serve(listener); err != nil {
		panicWithError(err, "failed to serve")
	}
}

func loadConfigOrPanic(cmd *cobra.Command) *Config {
	config, err := LoadConfig(cmd)
	if err != nil {
		panicWithError(err, "Failed to load configurations.")
	}
	return config
}

func panicWithError(err error, format string, args ...interface{}) {
	log.Fatalf("ERR: %v\n%s", err, fmt.Sprintf(format, args...))
}
