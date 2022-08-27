package storage

import "errors"

type logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

var (
	ErrConnFailed    = errors.New("failed to connect to database")
	ErrLinkExists    = errors.New("relation already exists")
	ErrOperationFail = errors.New("whoops, something went wrong during database operation")
	ErrEmptyResult   = errors.New("whoops, something went wrong, we have no idea what to show")
)
