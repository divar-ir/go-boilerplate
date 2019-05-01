package main

import (
	"log"
	"net"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/postview"

	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/cache"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/core"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/internal/provider"
	"google.golang.org/grpc"
)

func main() {
	providerInstance := provider.NewMemory()
	cacheInstance := cache.NewMemory()
	coreLogic := core.New(providerInstance, cacheInstance)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	postview.RegisterPostViewServer(s)
	pb.RegisterGreeterServer(s, coreLogic)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
