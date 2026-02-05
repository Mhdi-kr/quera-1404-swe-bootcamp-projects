package service

import "errors"

var ErrUserAlreadyRegistered = errors.New("user already registered")
var ErrUserNotFound = errors.New("user not found")
var ErrWrongCredentials = errors.New("wrong credentials")
