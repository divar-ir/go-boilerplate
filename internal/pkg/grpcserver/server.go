package grpcserver

import (
	"fmt"
	"git.cafebazaar.ir/bardia/lazyapi/pkg/appdetail"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
}

func New(appDetailServer appdetail.AppDetailServer, logger *logrus.Logger, listenPort int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		return nil, errors.Wrap(err, "failed to listen")
	}

	logEntry := logger.WithFields(map[string]interface{}{
		"app": "appdetaild",
	})

	interceptors := []grpc.UnaryServerInterceptor{
		grpclogrus.UnaryServerInterceptor(logEntry),
		grpcprometheus.UnaryServerInterceptor,
		grpcrecovery.UnaryServerInterceptor(),
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...)))
	appdetail.RegisterAppDetailServer(grpcServer, appDetailServer)

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
