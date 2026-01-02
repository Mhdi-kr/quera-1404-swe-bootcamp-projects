package repository

import "errors"

var ErrPostNotFound = errors.New("post does not exist")
var ErrUserNotFound = errors.New("user not found")
var ErrCommentNotFound = errors.New("comment not found")
