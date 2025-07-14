package converter

import (
	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
	"github.com/escoutdoor/vegetable_store/user_service/internal/utils/hasher"
	"github.com/gojaguar/jaguar/strings"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
)

func UserToProtoUser(user entity.User) *userv1.User {
	return &userv1.User{
		Id:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
}

func UsersToProtoUsers(users []entity.User) []*userv1.User {
	list := make([]*userv1.User, 0, len(users))
	for _, u := range users {
		list = append(list, UserToProtoUser(u))
	}

	return list
}

func ProtoUpdateUserRequestToUserUpdateOperation(req *userv1.UpdateUserRequest) (dto.UserUpdateOperation, error) {
	update := &userv1.UserUpdateOperation{}
	out := dto.UserUpdateOperation{ID: req.Update.GetUserId()}

	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.GetUpdateMask(), strings.PascalCase)
	if err != nil {
		return dto.UserUpdateOperation{}, errwrap.Wrap("create mask from the field mask", err)
	}

	if err := fieldmask_utils.StructToStruct(mask, req.Update, update); err != nil {
		return dto.UserUpdateOperation{}, errwrap.Wrap("copy struct to struct", err)
	}

	for _, p := range req.GetUpdateMask().GetPaths() {
		switch p {
		case "first_name":
			out.FirstName = &update.FirstName
		case "last_name":
			out.LastName = &update.LastName
		case "email":
			out.Email = &update.Email
		case "phone_number":
			out.PhoneNumber = &update.PhoneNumber
		case "password":
			*out.Password, err = hasher.HashPassword(update.Password)
			if err != nil {
				return dto.UserUpdateOperation{}, errwrap.Wrap("hash password", err)
			}
		}
	}

	return out, nil
}

func ProtoListUsersRequestToListUsersParams(req *userv1.ListUsersRequest) dto.ListUsersParams {
	return dto.ListUsersParams{
		Limit:   req.GetLimit(),
		Offset:  req.GetOffset(),
		UserIDs: req.GetUserIds(),
	}
}
