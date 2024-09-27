package data

import "errors"

var (
	ErrNotFound = errors.New("no rows matched for the query.")
)
