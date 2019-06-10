package errors

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		if err != nil {
			ctxlogrus.AddFields(ctx, Extras(err))
		}
	}()

	return handler(ctx, req)
}
