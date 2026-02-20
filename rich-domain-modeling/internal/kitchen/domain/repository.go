package domain

import (
	"context"
	"errors"
)

var (
	ErrVoucherNotFound      = errors.New("voucher not found")
	ErrVoucherAlreadyExists = errors.New("voucher already exists")
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderAlreadyExists   = errors.New("order already exists")
)

type VoucherRepository interface {
	SaveVoucher(context.Context, *CoffeeVoucher) error
	GetVoucherByID(context.Context, string) (*CoffeeVoucher, error)
	ListVouchersByCustomer(context.Context, string) ([]*CoffeeVoucher, error)
}

type OrderRepository interface {
	SaveOrder(context.Context, *CoffeeOrder) error
	GetOrderByID(context.Context, string) (*CoffeeOrder, error)
	ListOrdersByCustomer(context.Context, string) ([]*CoffeeOrder, error)
}
