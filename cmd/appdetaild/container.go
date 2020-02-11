//+build wireinject

package main

import (
	"context"
	"git.cafebazaar.ir/bardia/lazyapi/internal/app/core"
	"git.cafebazaar.ir/bardia/lazyapi/internal/pkg/grpcserver"
	"github.com/google/wire"
	"github.com/spf13/cobra"
)

func CreateServer(ctx context.Context, cmd *cobra.Command) (*grpcserver.Server, error) {
	panic(wire.Build(
		provideConfig,
		provideLogger,
		provideServer,
		core.New,
	))
}
