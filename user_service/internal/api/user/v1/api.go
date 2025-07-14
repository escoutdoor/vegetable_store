package v1

import (
	"fmt"

	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/service"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	firstNameFieldName   = "first_name"
	lastNameFieldName    = "last_name"
	emailFieldName       = "email"
	phoneNumberFieldName = "phone_number"
	passwordFieldName    = "password"

	minFirstNameLength = 1
	maxFirstNameLength = 20

	minLastNameLength = 1
	maxLastNameLength = 20
)

type Implementation struct {
	userService service.UserService
	userv1.UnimplementedUserServiceServer
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}

func validateUpdateUserRequest(req *userv1.UpdateUserRequest) error {
	if req.UpdateMask == nil {
		return fmt.Errorf("update_mask: unspecified")
	}
	if len(req.UpdateMask.Paths) == 0 {
		return fmt.Errorf("update_mask: paths: unspecified")
	}

	em := make(map[string]error, len(req.UpdateMask.Paths))
	for _, p := range req.UpdateMask.Paths {
		switch p {
		case firstNameFieldName:
			if err := validateFirstName(req.Update.FirstName); err != nil {
				em[firstNameFieldName] = err
			}
		case lastNameFieldName:
			if err := validateLastName(req.Update.LastName); err != nil {
				em[lastNameFieldName] = err
			}
		case emailFieldName:
			// TODO: validate
		case phoneNumberFieldName:
			// TODO: validate
		case passwordFieldName:
			// TODO: validate
		}
		if len(em) > 0 {
			vs := make([]*errdetails.BadRequest_FieldViolation, 0, len(em))
			for f, e := range em {
				vs = append(vs, &errdetails.BadRequest_FieldViolation{
					Field:       f,
					Description: e.Error(),
				})
			}

			st := status.New(codes.InvalidArgument, codes.InvalidArgument.String())

			badReq := &errdetails.BadRequest{FieldViolations: vs}

			detailedStatus, err := st.WithDetails(badReq)
			if err != nil {
				return err
			}

			return detailedStatus.Err()
		}
	}

	return nil
}

func validateFirstName(firstName string) error {
	if len(firstName) < minFirstNameLength {
		return fmt.Errorf("value length must be at least %d characters", minFirstNameLength)
	}

	if len(firstName) > maxFirstNameLength {
		return fmt.Errorf("value length must be at most %d characters", maxFirstNameLength)
	}

	return nil
}

func validateLastName(lastName string) error {
	if len(lastName) < minLastNameLength {
		return fmt.Errorf("value length must be at least %d characters", minLastNameLength)
	}

	if len(lastName) > maxLastNameLength {
		return fmt.Errorf("value length must be at most %d characters", maxLastNameLength)
	}

	return nil
}
