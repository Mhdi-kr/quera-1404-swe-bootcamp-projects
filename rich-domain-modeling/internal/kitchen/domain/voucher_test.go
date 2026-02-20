package domain_test

import (
	"testing"
	"time"

	kitchendomain "richdomainmodeling/internal/kitchen/domain"
)

func TestVoucherRedeemLifecycle(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	voucher, err := kitchendomain.NewCoffeeVoucher("voucher-1", "member-42", "vip seat", now)
	if err != nil {
		t.Fatalf("new voucher: %v", err)
	}

	if err := voucher.Redeem(now.Add(10 * time.Minute)); err != nil {
		t.Fatalf("redeem voucher: %v", err)
	}
	if voucher.Status() != kitchendomain.VoucherStatusRedeemed {
		t.Fatalf("expected redeemed status, got %s", voucher.Status())
	}
	if voucher.RedeemedAt() == nil {
		t.Fatalf("expected redeemed timestamp")
	}
}
