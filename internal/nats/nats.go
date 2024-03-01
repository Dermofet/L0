// Package nats provides functionality for interacting with NATS messaging system.
package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"L0/internal/cache"
	"L0/internal/entity"
	"L0/internal/repository"
)

// natsService represents a service for handling NATS messaging.
type natsService struct {
	logger          *zap.Logger
	orderRepository repository.OrderRepository
	cache           *cache.Cache
}

// NewNatsService creates a new instance of natsService.
func NewNatsService(orderRepository repository.OrderRepository, logger *zap.Logger, cache *cache.Cache) *natsService {
	return &natsService{
		logger:          logger,
		orderRepository: orderRepository,
		cache:           cache,
	}
}

// Subscribe subscribes to a NATS subject and processes incoming messages.
func (ns *natsService) Subscribe(ctx context.Context, clusterID string, clientID string, subject string, url string) error {
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return fmt.Errorf("can't connect to NATS: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	sub, err := conn.Subscribe(subject, ns.process)
	if err != nil {
		return fmt.Errorf("can't subscribe to NATS: %w", err)
	}

	<-ctx.Done()

	sub.Unsubscribe()

	return nil
}

// process handles incoming NATS messages.
func (ns *natsService) process(msg *stan.Msg) {
	var order entity.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		ns.logger.Error("can't decode message", zap.Error(err), zap.Any("data", msg.Data))
		return
	}

	id, err := ns.orderRepository.Create(context.Background(), &order)
	if err != nil {
		ns.logger.Error("can't save order to database", zap.Error(err))
		return
	}

	ns.cache.Set(id, &order)
}
