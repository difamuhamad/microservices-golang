package error

import "errors"

var (
	ErrFieldScheduleNotFound = errors.New("field schedule not found")
	ErrFieldScheduleNotExist = errors.New("field schedule already exist")
)

var ErrFieldScheduleErrors = []error{
	ErrFieldScheduleNotFound,
	ErrFieldScheduleNotExist,
}
