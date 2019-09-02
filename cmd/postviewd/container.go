//+build wireinject

package main

import (
	"context"
	"github.com/cafebazaar/go-boilerplate/internal/app/core"
	"github.com/cafebazaar/go-boilerplate/internal/pkg/grpcserver"
	"github.com/cafebazaar/go-boilerplate/internal/app/provider"
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
