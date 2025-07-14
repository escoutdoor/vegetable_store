package interceptor

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/escoutdoor/vegetable_store/common/pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func ValidationUnaryServerInterceptor(validator protovalidate.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("unsupported message type: %T", msg)
		}

		if err := validator.Validate(msg); err != nil {
			return nil, grpcutil.ProtoValidationError(err)
		}

		return handler(ctx, req)
	}
}
