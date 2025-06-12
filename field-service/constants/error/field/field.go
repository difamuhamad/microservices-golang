package error

import "errors"

var (
	ErrFieldNotFound         = errors.New("field not found")
	ErrFieldScheduleNotFound = errors.New("field schedule not found")
)

var FieldErrors = []error{
	ErrFieldNotFound,
	ErrFieldScheduleNotFound,
}
