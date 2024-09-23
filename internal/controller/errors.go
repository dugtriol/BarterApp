package controller

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotCreated       = errors.New("not created")
	ErrNotValid         = errors.New("failed to unmarshall")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrInput            = errors.New("input errors")
)
