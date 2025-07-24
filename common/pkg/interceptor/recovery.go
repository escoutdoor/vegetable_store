package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/escoutdoor/vegetable_store/common/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoverUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if v := recover(); v != nil {
				logger.ErrorKV(ctx, "recover panic",
					"panic", v,
					"stacktrace", string(debug.Stack()),
					"operation", info.FullMethod,
					"component", "interceptor",
				)

				err = status.Error(codes.Internal, "internal server error") // return error
			}
		}()

		return handler(ctx, req)
	}
}
