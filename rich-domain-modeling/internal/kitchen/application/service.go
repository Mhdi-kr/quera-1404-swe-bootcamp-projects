package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	kitchendomain "richdomainmodeling/internal/kitchen/domain"
	sharedapp "richdomainmodeling/internal/shared/application"
	shareddomain "richdomainmodeling/internal/shared/domain"
	theaterdomain "richdomainmodeling/internal/theater/domain"
)

type Service struct {
	voucherRepository kitchendomain.VoucherRepository
	orderRepository   kitchendomain.OrderRepository
	eventBus          sharedapp.EventBus
	tracer            trace.Tracer
}

func NewService(
	voucherRepository kitchendomain.VoucherRepository,
	orderRepository kitchendomain.OrderRepository,
	eventBus sharedapp.EventBus,
) *Service {
	return &Service{
		voucherRepository: voucherRepository,
		orderRepository:   orderRepository,
		eventBus:          eventBus,
		tracer:            otel.Tracer("richdomainmodeling/internal/kitchen/application"),
	}
}

func (s *Service) IssueComplimentaryCoffee(
	ctx context.Context,
	voucherID string,
	customerID string,
	source string,
	issuedAt time.Time,
) (*kitchendomain.CoffeeVoucher, error) {
	_, domainSpan := s.tracer.Start(ctx, "domain.kitchen.voucher.new")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("voucher.id", voucherID),
		attribute.String("customer.id", customerID),
	)
	voucher, err := kitchendomain.NewCoffeeVoucher(voucherID, customerID, source, issuedAt)
	if err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return nil, err
	}
	domainSpan.End()

	if err := s.voucherRepository.SaveVoucher(ctx, voucher); err != nil {
		return nil, err
	}
	if err := s.publish(ctx, voucher.PullEvents()); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *Service) PlacePaidOrder(
	ctx context.Context,
	orderID string,
	customerID string,
	drink string,
	priceCents int,
	orderedAt time.Time,
) (*kitchendomain.CoffeeOrder, error) {
	_, domainSpan := s.tracer.Start(ctx, "domain.kitchen.order.new_paid")
	domainSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("order.id", orderID),
		attribute.String("customer.id", customerID),
	)
	order, err := kitchendomain.NewPaidCoffeeOrder(orderID, customerID, drink, priceCents, orderedAt)
	if err != nil {
		domainSpan.RecordError(err)
		domainSpan.SetStatus(codes.Error, err.Error())
		domainSpan.End()
		return nil, err
	}
	domainSpan.End()

	if err := s.orderRepository.SaveOrder(ctx, order); err != nil {
		return nil, err
	}
	if err := s.publish(ctx, order.PullEvents()); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) RedeemVoucher(
	ctx context.Context,
	voucherID string,
	orderID string,
	drink string,
	redeemedAt time.Time,
) (*kitchendomain.CoffeeOrder, error) {
	voucher, err := s.voucherRepository.GetVoucherByID(ctx, voucherID)
	if err != nil {
		return nil, err
	}

	_, redeemSpan := s.tracer.Start(ctx, "domain.kitchen.voucher.redeem")
	redeemSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("voucher.id", voucherID),
	)
	if err := voucher.Redeem(redeemedAt); err != nil {
		redeemSpan.RecordError(err)
		redeemSpan.SetStatus(codes.Error, err.Error())
		redeemSpan.End()
		return nil, err
	}
	redeemSpan.End()

	_, orderSpan := s.tracer.Start(ctx, "domain.kitchen.order.new_complimentary")
	orderSpan.SetAttributes(
		attribute.String("layer", "domain"),
		attribute.String("order.id", orderID),
	)
	order, err := kitchendomain.NewComplimentaryCoffeeOrder(orderID, voucher.CustomerID(), drink, redeemedAt)
	if err != nil {
		orderSpan.RecordError(err)
		orderSpan.SetStatus(codes.Error, err.Error())
		orderSpan.End()
		return nil, err
	}
	orderSpan.End()

	if err := s.voucherRepository.SaveVoucher(ctx, voucher); err != nil {
		return nil, err
	}
	if err := s.orderRepository.SaveOrder(ctx, order); err != nil {
		return nil, err
	}

	events := append(voucher.PullEvents(), order.PullEvents()...)
	if err := s.publish(ctx, events); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) ListVouchersByCustomer(ctx context.Context, customerID string) ([]*kitchendomain.CoffeeVoucher, error) {
	return s.voucherRepository.ListVouchersByCustomer(ctx, customerID)
}

func (s *Service) HandleVIPSeatPurchased(ctx context.Context, event shareddomain.DomainEvent) error {
	vip, ok := event.(theaterdomain.VIPSeatPurchased)
	if !ok {
		return nil
	}

	voucherID := fmt.Sprintf("free-coffee-%s", vip.TicketID)
	source := fmt.Sprintf("VIP seat %s for show %s", vip.SeatNumber, vip.ShowID)
	_, err := s.IssueComplimentaryCoffee(ctx, voucherID, vip.CustomerID, source, vip.OccurredAt())
	if err != nil && !errors.Is(err, kitchendomain.ErrVoucherAlreadyExists) {
		return err
	}

	return nil
}

func (s *Service) publish(ctx context.Context, events []shareddomain.DomainEvent) error {
	if s.eventBus == nil || len(events) == 0 {
		return nil
	}
	return s.eventBus.Publish(ctx, events...)
}
