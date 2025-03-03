package user

import "errors"

var (
	ErrStorageFail   = errors.New("storage error")
	ErrWrongPassword = errors.New("wrong password")
)
