package domain

import (
	"errors"
	"time"

	shareddomain "richdomainmodeling/internal/shared/domain"
)

var (
	ErrInvalidOrderID    = errors.New("order id is required")
	ErrInvalidDrink      = errors.New("drink is required")
	ErrInvalidOrderPrice = errors.New("order price must be positive")
	ErrInvalidOrderOwner = errors.New("order customer id is required")
)

type CoffeeOrder struct {
	recorder      shareddomain.EventRecorder
	id            string
	customerID    string
	drink         string
	priceCents    int
	complimentary bool
	createdAt     time.Time
}

func NewPaidCoffeeOrder(id, customerID, drink string, priceCents int, createdAt time.Time) (*CoffeeOrder, error) {
	switch {
	case id == "":
		return nil, ErrInvalidOrderID
	case customerID == "":
		return nil, ErrInvalidOrderOwner
	case drink == "":
		return nil, ErrInvalidDrink
	case priceCents <= 0:
		return nil, ErrInvalidOrderPrice
	}

	order := &CoffeeOrder{
		id:         id,
		customerID: customerID,
		drink:      drink,
		priceCents: priceCents,
		createdAt:  createdAt,
	}
	order.recorder.Record(NewPaidCoffeeOrdered(id, customerID, drink, priceCents, createdAt))

	return order, nil
}

func NewComplimentaryCoffeeOrder(id, customerID, drink string, createdAt time.Time) (*CoffeeOrder, error) {
	switch {
	case id == "":
		return nil, ErrInvalidOrderID
	case customerID == "":
		return nil, ErrInvalidOrderOwner
	case drink == "":
		return nil, ErrInvalidDrink
	}

	return &CoffeeOrder{
		id:            id,
		customerID:    customerID,
		drink:         drink,
		priceCents:    0,
		complimentary: true,
		createdAt:     createdAt,
	}, nil
}

func (o *CoffeeOrder) ID() string {
	return o.id
}

func (o *CoffeeOrder) CustomerID() string {
	return o.customerID
}

func (o *CoffeeOrder) Drink() string {
	return o.drink
}

func (o *CoffeeOrder) PriceCents() int {
	return o.priceCents
}

func (o *CoffeeOrder) Complimentary() bool {
	return o.complimentary
}

func (o *CoffeeOrder) CreatedAt() time.Time {
	return o.createdAt
}

func (o *CoffeeOrder) PullEvents() []shareddomain.DomainEvent {
	return o.recorder.PullEvents()
}
