package inmemory

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	kitchendomain "richdomainmodeling/internal/kitchen/domain"
)

type VoucherRepository struct {
	mu       sync.RWMutex
	vouchers map[string]*kitchendomain.CoffeeVoucher
	tracer   trace.Tracer
}

func NewVoucherRepository() *VoucherRepository {
	return &VoucherRepository{
		vouchers: make(map[string]*kitchendomain.CoffeeVoucher),
		tracer:   otel.Tracer("richdomainmodeling/internal/infrastructure/persistence/inmemory"),
	}
}

func (r *VoucherRepository) SaveVoucher(ctx context.Context, voucher *kitchendomain.CoffeeVoucher) (err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.voucher.save")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("voucher.id", voucher.ID()),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.vouchers[voucher.ID()]; ok && existing != voucher {
		err = kitchendomain.ErrVoucherAlreadyExists
		return err
	}

	r.vouchers[voucher.ID()] = voucher
	return nil
}

func (r *VoucherRepository) GetVoucherByID(ctx context.Context, voucherID string) (voucher *kitchendomain.CoffeeVoucher, err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.voucher.get_by_id")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("voucher.id", voucherID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	voucher, ok := r.vouchers[voucherID]
	if !ok {
		err = kitchendomain.ErrVoucherNotFound
		return nil, err
	}
	return voucher, nil
}

func (r *VoucherRepository) ListVouchersByCustomer(ctx context.Context, customerID string) (vouchers []*kitchendomain.CoffeeVoucher, err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.voucher.list_by_customer")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("customer.id", customerID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.SetAttributes(attribute.Int("voucher.count", len(vouchers)))
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	vouchers = make([]*kitchendomain.CoffeeVoucher, 0)
	for _, voucher := range r.vouchers {
		if voucher.CustomerID() == customerID {
			vouchers = append(vouchers, voucher)
		}
	}
	return vouchers, nil
}

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*kitchendomain.CoffeeOrder
	tracer trace.Tracer
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]*kitchendomain.CoffeeOrder),
		tracer: otel.Tracer("richdomainmodeling/internal/infrastructure/persistence/inmemory"),
	}
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *kitchendomain.CoffeeOrder) (err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.order.save")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("order.id", order.ID()),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.orders[order.ID()]; ok && existing != order {
		err = kitchendomain.ErrOrderAlreadyExists
		return err
	}

	r.orders[order.ID()] = order
	return nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (order *kitchendomain.CoffeeOrder, err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.order.get_by_id")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("order.id", orderID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderID]
	if !ok {
		err = kitchendomain.ErrOrderNotFound
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) ListOrdersByCustomer(ctx context.Context, customerID string) (orders []*kitchendomain.CoffeeOrder, err error) {
	_, span := r.tracer.Start(ctx, "repo.kitchen.order.list_by_customer")
	span.SetAttributes(
		attribute.String("layer", "repo"),
		attribute.String("customer.id", customerID),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.SetAttributes(attribute.Int("order.count", len(orders)))
		span.End()
	}()

	r.mu.RLock()
	defer r.mu.RUnlock()

	orders = make([]*kitchendomain.CoffeeOrder, 0)
	for _, order := range r.orders {
		if order.CustomerID() == customerID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
