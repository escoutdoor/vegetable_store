package interceptor

import (
	"context"
	"errors"

	apperrors "github.com/escoutdoor/vegetable_store/order_service/internal/errors"
	"github.com/escoutdoor/vegetable_store/order_service/internal/errors/codes"
	"google.golang.org/grpc"

	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)

		if _, ok := status.FromError(err); ok {
			return resp, err
		}

		appErr := new(apperrors.Error)
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case codes.OrderNotFound:
				err = status.Error(grpccodes.NotFound, appErr.Error())
			case codes.InsufficientInventory:
				err = status.Error(grpccodes.InvalidArgument, appErr.Error())
			case codes.VegetablesNotFound:
				err = status.Error(grpccodes.NotFound, appErr.Error())
			}
		} else {
			err = status.Error(grpccodes.Internal, "internal server error")
		}

		return resp, err
	}
}
