package repo

import "errors"

var (
	ErrNoRows      = errors.New("no rows in result set")
	ErrTooManyRows = errors.New("too many rows in result set")
)
