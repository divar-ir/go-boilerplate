package grpcserver

import (
	"fmt"
	"github.com/cafebazaar/go-boilerplate/pkg/errors"
	"github.com/cafebazaar/go-boilerplate/pkg/postview"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
}

func New(postViewServer postview.PostViewServer, logger *logrus.Logger, listenPort int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	logEntry := logger.WithFields(map[string]interface{}{
		"app": "postviewd",
	})

	interceptors := []grpc.UnaryServerInterceptor{
		grpclogrus.UnaryServerInterceptor(logEntry),
		errors.UnaryServerInterceptor,
		grpcprometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	postview.RegisterPostViewServer(grpcServer, postViewServer)

	return &Server{
		listener: listener,
		server:   grpcServer,
	}, nil
}

func (s *Server) Serve() error {
	if err := s.server.Serve(s.listener); err != nil {
		return errors.Wrap(err, "failed to serve")
	}
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
