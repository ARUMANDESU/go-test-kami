package storage

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrResevationConflict = errors.New("reservation conflict")
)
