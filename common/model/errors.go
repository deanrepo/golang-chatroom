package model

import "errors"

var (
	// ErrUserNotExists represents user not exists
	ErrUserNotExists = errors.New("user not exists")
	// ErrUserExists represents user exists
	ErrUserExists = errors.New("user exists")
	// ErrWrongPwd represents wrong password
	ErrWrongPwd = errors.New("wrong password")
)
