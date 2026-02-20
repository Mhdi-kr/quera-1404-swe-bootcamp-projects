package integration_test

import (
	"context"
	"testing"
	"time"

	"richdomainmodeling/internal/infrastructure/eventbus"
	"richdomainmodeling/internal/infrastructure/persistence/inmemory"
	kitchenapp "richdomainmodeling/internal/kitchen/application"
	kitchendomain "richdomainmodeling/internal/kitchen/domain"
	theaterapp "richdomainmodeling/internal/theater/application"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

func TestVIPSeatPurchaseIssuesFreeCoffeeVoucher(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bus := eventbus.NewInMemoryEventBus()
	showRepo := inmemory.NewShowRepository()
	voucherRepo := inmemory.NewVoucherRepository()
	orderRepo := inmemory.NewOrderRepository()

	theaterService := theaterapp.NewService(showRepo, bus)
	kitchenService := kitchenapp.NewService(voucherRepo, orderRepo, bus)

	bus.Subscribe(theaterdomain.EventVIPSeatPurchased, kitchenService.HandleVIPSeatPurchased)

	startsAt := time.Now().UTC().Add(24 * time.Hour)
	_, err := theaterService.ScheduleShow(ctx, "show-vip-test", "VIP Night", startsAt, []theaterdomain.Seat{
		{Number: "V1", Tier: theaterdomain.SeatTierVIP, PriceCents: 9000},
	})
	if err != nil {
		t.Fatalf("schedule show: %v", err)
	}

	ticket, err := theaterService.PurchaseSeat(ctx, "show-vip-test", "member-1", "V1", startsAt.Add(-time.Hour))
	if err != nil {
		t.Fatalf("purchase seat: %v", err)
	}
	if !ticket.IncludesFreeCoffee {
		t.Fatalf("expected VIP ticket to include free coffee")
	}

	vouchers, err := kitchenService.ListVouchersByCustomer(ctx, "member-1")
	if err != nil {
		t.Fatalf("list vouchers: %v", err)
	}
	if len(vouchers) != 1 {
		t.Fatalf("expected exactly one free coffee voucher, got %d", len(vouchers))
	}
	if vouchers[0].Status() != kitchendomain.VoucherStatusIssued {
		t.Fatalf("expected issued voucher status, got %s", vouchers[0].Status())
	}
}
