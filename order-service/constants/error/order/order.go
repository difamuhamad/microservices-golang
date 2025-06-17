package error

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrFiledAlreadyBooked = errors.New("filed schedule already booked")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrFiledAlreadyBooked,
}
