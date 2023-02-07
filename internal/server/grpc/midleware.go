package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func loggingMiddleware(log Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		response, err := handler(ctx, req)
		log.Info(fmt.Sprintf("%s %s %s %s", time.Now().Format(time.RFC822Z), info.FullMethod, req, time.Since(start)))
		return response, err
	}
}
