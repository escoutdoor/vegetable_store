package v1

import (
	"context"

	authv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/auth/v1"
)

func (i *Implementation) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	userID, err := i.authService.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return &authv1.ValidateTokenResponse{
		UserId: userID,
	}, nil
}
