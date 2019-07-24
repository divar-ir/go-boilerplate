//+build wireinject

package main

import (
	"context"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/app/core"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/grpcserver"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

func CreateServer(ctx context.Context, cmd *cobra.Command) (*grpcserver.Server, error) {
	panic(wire.Build(
		provideConfig,
		provideLogger,
		provideProvider,
		provideCache,
		providePrometheus,
		provideServer,
		core.New,
	))
}

func CreateProvider(ctx context.Context, cmd *cobra.Command) (provider.PostProvider, error) {
	panic(wire.Build(
		provideConfig,
		provideLogger,
		provideProvider,
		providePrometheus,
	))
}
