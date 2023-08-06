package errors

import "fmt"

var (
	TokenTypeError        = fmt.Errorf("token type unknown")
	TokenExtractUserError = fmt.Errorf("type assertion to user err")
)

const (
	NoError = iota
	InternalError
	GeneralError
)

const (
	HashPasswordError = iota + 1000
)

const (
	AuthServiceGeneralErr = iota + 2000
	AuthServiceWrongPasswordErr
	AuthServiceTokenGenerationErr
	AuthServiceUserNotVerified
	UserServiceCreateUserErr
	UserServiceUserAlreadyExists
)
