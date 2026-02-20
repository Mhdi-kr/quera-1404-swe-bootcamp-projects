package domain

import (
	"errors"
	"time"

	shareddomain "richdomainmodeling/internal/shared/domain"
)

type VoucherStatus string

const (
	VoucherStatusIssued   VoucherStatus = "issued"
	VoucherStatusRedeemed VoucherStatus = "redeemed"
)

var (
	ErrInvalidVoucherID     = errors.New("voucher id is required")
	ErrInvalidVoucherSource = errors.New("voucher source is required")
	ErrInvalidVoucherOwner  = errors.New("voucher customer id is required")
	ErrVoucherNotRedeemable = errors.New("voucher is not redeemable")
)

type CoffeeVoucher struct {
	recorder   shareddomain.EventRecorder
	id         string
	customerID string
	source     string
	status     VoucherStatus
	issuedAt   time.Time
	redeemedAt *time.Time
}

func NewCoffeeVoucher(id, customerID, source string, issuedAt time.Time) (*CoffeeVoucher, error) {
	switch {
	case id == "":
		return nil, ErrInvalidVoucherID
	case customerID == "":
		return nil, ErrInvalidVoucherOwner
	case source == "":
		return nil, ErrInvalidVoucherSource
	}

	voucher := &CoffeeVoucher{
		id:         id,
		customerID: customerID,
		source:     source,
		status:     VoucherStatusIssued,
		issuedAt:   issuedAt,
	}
	voucher.recorder.Record(NewComplimentaryCoffeeIssued(id, customerID, source, issuedAt))

	return voucher, nil
}

func (v *CoffeeVoucher) ID() string {
	return v.id
}

func (v *CoffeeVoucher) CustomerID() string {
	return v.customerID
}

func (v *CoffeeVoucher) Source() string {
	return v.source
}

func (v *CoffeeVoucher) Status() VoucherStatus {
	return v.status
}

func (v *CoffeeVoucher) IssuedAt() time.Time {
	return v.issuedAt
}

func (v *CoffeeVoucher) RedeemedAt() *time.Time {
	return v.redeemedAt
}

func (v *CoffeeVoucher) Redeem(at time.Time) error {
	if v.status != VoucherStatusIssued {
		return ErrVoucherNotRedeemable
	}

	v.status = VoucherStatusRedeemed
	v.redeemedAt = &at
	v.recorder.Record(NewCoffeeVoucherRedeemed(v.id, v.customerID, at))

	return nil
}

func (v *CoffeeVoucher) PullEvents() []shareddomain.DomainEvent {
	return v.recorder.PullEvents()
}
