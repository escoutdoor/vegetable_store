package converter

import (
	authv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/auth/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
)

func ProtoRegisterRequestToCreateUserParams(req *authv1.RegisterRequest) dto.CreateUserParams {
	return dto.CreateUserParams{
		FirstName:   req.GetFirstName(),
		LastName:    req.GetLastName(),
		Email:       req.GetEmail(),
		Password:    req.GetPassword(),
		PhoneNumber: req.GetPhoneNumber().String(),
	}
}
