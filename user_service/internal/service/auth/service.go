package auth

import (
	"github.com/escoutdoor/vegetable_store/user_service/internal/repository"
	"github.com/escoutdoor/vegetable_store/user_service/internal/utils/token"
)

type service struct {
	userRepository repository.UserRepository
	tokenProvider  token.Provider
}

func NewService(userRepository repository.UserRepository, tokenProvider token.Provider) *service {
	return &service{
		userRepository: userRepository,
		tokenProvider:  tokenProvider,
	}
}
