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
	cache *cache.Cache
}

// NewOrderInteractor creates a new instance of orderInteractor.
func NewOrderInteractor(repo repository.OrderRepository, cache *cache.Cache) *orderInteractor {
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
