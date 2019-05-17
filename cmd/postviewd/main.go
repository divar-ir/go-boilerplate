package main

import (
	"log"
	"net"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/app/core"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/pkg/provider"
	"google.golang.org/grpc"
)

func main() {
	providerInstance := provider.NewMemory()
	cacheInstance := cache.NewMemory()
	servicer := core.New(providerInstance, cacheInstance)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	postview.RegisterPostViewServer(grpcServer, servicer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
