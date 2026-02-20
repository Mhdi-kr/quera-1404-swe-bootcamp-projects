package domain

import (
	"context"
	"errors"
)

var ErrShowNotFound = errors.New("show not found")

type Repository interface {
	Save(context.Context, *Show) error
	GetByID(context.Context, string) (*Show, error)
}
