package codes

type Code string

const (
	UserNotFound          Code = "USER_NOT_FOUND"
	EmailAlreadyExists    Code = "EMAIL_ALREADY_EXISTS"
	IncorrectCreadentials Code = "INCORRECT_CREADENTIALS"
	JwtTokenExpired       Code = "JWT_TOKEN_EXPIRED"
	InvalidJwtToken       Code = "INVALID_JWT_TOKEN"
)
