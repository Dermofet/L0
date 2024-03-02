// Package usecase provides implementations for order use cases.
package usecase

import (
	"L0/internal/cache"
	"L0/internal/entity"
	"L0/internal/repository"
	"context"
	"fmt"
)

// orderInteractor implements the OrderInteractor interface.
type orderInteractor struct {
	repo  repository.OrderRepository
	cache cache.Cache
}

// NewOrderInteractor creates a new instance of orderInteractor.
func NewOrderInteractor(repo repository.OrderRepository, cache cache.Cache) *orderInteractor {
	return &orderInteractor{
		repo:  repo,
		cache: cache,
	}
}

// Create creates a new order.
func (u *orderInteractor) Create(ctx context.Context, order *entity.Order) error {
	id, err := u.repo.Create(ctx, order)
	if err != nil {
		return fmt.Errorf("can't create order by repository: %w", err)
	}

	u.cache.Set(id, order)

	return nil
}

// GetByUid retrieves an order by its UID.
func (u *orderInteractor) GetByUid(ctx context.Context, uid string) (*entity.Order, error) {
	orderCache, ok := u.cache.Get(uid)
	if ok {
		return orderCache.(*entity.Order), nil
	}

	order, err := u.repo.GetByUid(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("can't get order by uid from repository: %w", err)
	}

	u.cache.Set(uid, order)

	return order, nil
}

// GetAll retrieves all orders.
func (u *orderInteractor) GetAll(ctx context.Context) ([]*entity.Order, error) {
	ordersCache, ok := u.cache.GetAll()
	if ok {
		var orders []*entity.Order
		for _, order := range ordersCache {
			orders = append(orders, order.(*entity.Order))
		}
		return orders, nil
	}

	orders, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get all orders from repository: %w", err)
	}

	return orders, nil
}

// Delete deletes an order.
func (u *orderInteractor) Delete(ctx context.Context, uid string) error {
	err := u.repo.Delete(ctx, uid)
	if err != nil {
		return fmt.Errorf("can't delete order: %w", err)
	}

	u.cache.Delete(uid)

	return nil
}
