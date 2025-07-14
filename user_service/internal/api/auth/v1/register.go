package v1

import (
	"context"

	authv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/auth/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/converter"
)

func (i *Implementation) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	userID, err := i.authService.Register(ctx, converter.ProtoRegisterRequestToCreateUserParams(req))
	if err != nil {
		return nil, err
	}

	return &authv1.RegisterResponse{UserId: userID}, nil
}
