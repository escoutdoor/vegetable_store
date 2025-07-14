package v1

import (
	authv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/auth/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/service"
)

type Implementation struct {
	authService service.AuthService
	authv1.UnimplementedAuthServiceServer
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
