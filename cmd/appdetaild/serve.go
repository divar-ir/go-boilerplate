package main

import (
	"context"
	"git.cafebazaar.ir/bardia/lazyapi/internal/pkg/grpcserver"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"git.cafebazaar.ir/bardia/lazyapi/pkg/appdetail"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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

func provideServer(server appdetail.AppDetailServer, config *Config, logger *logrus.Logger) (*grpcserver.Server, error) {
	return grpcserver.New(server, logger, config.ListenPort)
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
	if err == nil {
		return
	}

	logrus.WithError(err).Panicf(format, args...)
}
