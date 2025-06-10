package error

import "errors"

const (
	Success = "success"
	Error   = "error"
)

var (
	InternalServerError = errors.New("Internal Server Error")
	ErrSqlError         = errors.New("database server failed to execute query")
	ErrTooManyRequest   = errors.New("too many request")
	ErrUnauthorized     = errors.New("Unauthorized")
	ErrInvalidToken     = errors.New("invalid token")
	ErrForbidden        = errors.New("forbidden")
)

var GeneralErrors = []error{
	InternalServerError,
	ErrSqlError,
	ErrTooManyRequest,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrForbidden,
}
