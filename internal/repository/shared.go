package repository

import "errors"

var (
	ErrLuhnNumberIsInvalid = errors.New("invalid luhn number")
)
