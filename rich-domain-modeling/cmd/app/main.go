package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"richdomainmodeling/internal/infrastructure/eventbus"
	"richdomainmodeling/internal/infrastructure/persistence/inmemory"
	kitchenapp "richdomainmodeling/internal/kitchen/application"
	kitchendomain "richdomainmodeling/internal/kitchen/domain"
	reservationapp "richdomainmodeling/internal/reservation/application"
	theaterapp "richdomainmodeling/internal/theater/application"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

func main() {
	ctx := context.Background()
	bus := eventbus.NewInMemoryEventBus()

	reservationRepo := inmemory.NewReservationRepository()
	showRepo := inmemory.NewShowRepository()
	voucherRepo := inmemory.NewVoucherRepository()
	orderRepo := inmemory.NewOrderRepository()

	reservationService := reservationapp.NewService(reservationRepo, bus)
	theaterService := theaterapp.NewService(showRepo, bus)
	kitchenService := kitchenapp.NewService(voucherRepo, orderRepo, bus)

	bus.Subscribe(theaterdomain.EventVIPSeatPurchased, kitchenService.HandleVIPSeatPurchased)

	now := time.Now().UTC()

	// Reservation domain usual affairs.
	reservation, err := reservationService.ReserveWorkspace(
		ctx,
		"res-1001",
		"workspace-42",
		"customer-7",
		now.Add(48*time.Hour),
		now.Add(56*time.Hour),
		now,
	)
	if err != nil {
		log.Fatalf("reserve workspace: %v", err)
	}
	if err := reservationService.ConfirmReservation(ctx, reservation.ID(), now.Add(1*time.Hour)); err != nil {
		log.Fatalf("confirm reservation: %v", err)
	}

	// Theater domain usual affairs.
	_, err = theaterService.ScheduleShow(
		ctx,
		"show-1",
		"Foundations of DDD",
		now.Add(72*time.Hour),
		[]theaterdomain.Seat{
			{Number: "A1", Tier: theaterdomain.SeatTierStandard, PriceCents: 2500},
			{Number: "A2", Tier: theaterdomain.SeatTierVIP, PriceCents: 5000},
		},
	)
	if err != nil {
		log.Fatalf("schedule show: %v", err)
	}

	if _, err := theaterService.PurchaseSeat(ctx, "show-1", "customer-7", "A1", now.Add(2*time.Hour)); err != nil {
		log.Fatalf("purchase standard seat: %v", err)
	}

	vipTicket, err := theaterService.PurchaseSeat(ctx, "show-1", "customer-7", "A2", now.Add(3*time.Hour))
	if err != nil {
		log.Fatalf("purchase vip seat: %v", err)
	}

	// Kitchen usual affairs: paid order and voucher redemption.
	if _, err := kitchenService.PlacePaidOrder(ctx, "order-paid-1", "customer-7", "espresso", 450, now.Add(4*time.Hour)); err != nil {
		log.Fatalf("place paid order: %v", err)
	}

	vouchers, err := kitchenService.ListVouchersByCustomer(ctx, "customer-7")
	if err != nil {
		log.Fatalf("list vouchers: %v", err)
	}

	for _, voucher := range vouchers {
		if voucher.Status() == kitchendomain.VoucherStatusIssued {
			if _, err := kitchenService.RedeemVoucher(
				ctx,
				voucher.ID(),
				"order-free-1",
				"cappuccino",
				now.Add(5*time.Hour),
			); err != nil {
				log.Fatalf("redeem voucher %s: %v", voucher.ID(), err)
			}
		}
	}

	fmt.Printf("Reservation %s status: %s\n", reservation.ID(), reservation.Status())
	fmt.Printf("VIP ticket %s includes free coffee: %t\n", vipTicket.ID, vipTicket.IncludesFreeCoffee)
	fmt.Printf("Issued vouchers for customer-7: %d\n", len(vouchers))
}
