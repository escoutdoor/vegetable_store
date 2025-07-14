package dto

type CreateUserParams struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Password    string
}

type UserUpdateOperation struct {
	ID          string
	FirstName   *string
	LastName    *string
	Email       *string
	PhoneNumber *string
	Password    *string
}

type ListUsersParams struct {
	Limit   int64
	Offset  int64
	UserIDs []string
}
