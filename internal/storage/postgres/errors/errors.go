package errors

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

var (
	ErrAdminExists   = errors.New("admin already exists")
	ErrAdminNotFound = errors.New("admin not found")
)
