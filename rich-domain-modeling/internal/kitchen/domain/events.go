package domain

import "time"

const (
	EventComplimentaryCoffeeIssued = "kitchen.complimentary.coffee.issued"
	EventCoffeeVoucherRedeemed     = "kitchen.coffee.voucher.redeemed"
	EventPaidCoffeeOrdered         = "kitchen.coffee.paid.ordered"
)

type ComplimentaryCoffeeIssued struct {
	VoucherID  string
	CustomerID string
	Source     string
	OccurredOn time.Time
}

func NewComplimentaryCoffeeIssued(voucherID, customerID, source string, occurredOn time.Time) ComplimentaryCoffeeIssued {
	return ComplimentaryCoffeeIssued{
		VoucherID:  voucherID,
		CustomerID: customerID,
		Source:     source,
		OccurredOn: occurredOn,
	}
}

func (e ComplimentaryCoffeeIssued) EventName() string {
	return EventComplimentaryCoffeeIssued
}

func (e ComplimentaryCoffeeIssued) OccurredAt() time.Time {
	return e.OccurredOn
}

type CoffeeVoucherRedeemed struct {
	VoucherID  string
	CustomerID string
	OccurredOn time.Time
}

func NewCoffeeVoucherRedeemed(voucherID, customerID string, occurredOn time.Time) CoffeeVoucherRedeemed {
	return CoffeeVoucherRedeemed{
		VoucherID:  voucherID,
		CustomerID: customerID,
		OccurredOn: occurredOn,
	}
}

func (e CoffeeVoucherRedeemed) EventName() string {
	return EventCoffeeVoucherRedeemed
}

func (e CoffeeVoucherRedeemed) OccurredAt() time.Time {
	return e.OccurredOn
}

type PaidCoffeeOrdered struct {
	OrderID    string
	CustomerID string
	Drink      string
	PriceCents int
	OccurredOn time.Time
}

func NewPaidCoffeeOrdered(orderID, customerID, drink string, priceCents int, occurredOn time.Time) PaidCoffeeOrdered {
	return PaidCoffeeOrdered{
		OrderID:    orderID,
		CustomerID: customerID,
		Drink:      drink,
		PriceCents: priceCents,
		OccurredOn: occurredOn,
	}
}

func (e PaidCoffeeOrdered) EventName() string {
	return EventPaidCoffeeOrdered
}

func (e PaidCoffeeOrdered) OccurredAt() time.Time {
	return e.OccurredOn
}
