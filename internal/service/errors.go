package service

import "fmt"

var (
	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrCannotCreateUser  = fmt.Errorf("cannot create user")
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrCannotGetUser     = fmt.Errorf("cannot get user")
	ErrCannotUpdateUser  = fmt.Errorf("cannot update user")

	ErrProductAlreadyExists = fmt.Errorf("product already exists")
	ErrCannotCreateProduct  = fmt.Errorf("cannot create product")
	ErrCannotGetProduct     = fmt.Errorf("cannot get product")
	ErrCannotUpdateProduct  = fmt.Errorf("cannot update product")

	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrCannotCreate  = fmt.Errorf("cannot create")
	ErrCannotGet     = fmt.Errorf("cannot get")
)
