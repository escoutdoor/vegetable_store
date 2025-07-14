package user

import "github.com/escoutdoor/vegetable_store/user_service/internal/repository"

type service struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
