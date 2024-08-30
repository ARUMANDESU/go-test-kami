package reservation

import "errors"

var (
	ErrInternal = errors.New("internal error")

	ErrReservationConflict = errors.New("reservation conflict")
	ErrInvalidArgument     = errors.New("invalid argument")
)
