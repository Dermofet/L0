// Package nats provides functionality for interacting with NATS messaging system.
package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"

	"L0/internal/cache"
	"L0/internal/entity"
	"L0/internal/repository"
)

// natsService represents a service for handling NATS messaging.
type natsService struct {
	orderRepository repository.OrderRepository
	cache           cache.Cache
	connect         stan.Conn
	subject         string
}

// NewNatsService creates a new instance of natsService.
func NewNatsService(orderRepository repository.OrderRepository, cache cache.Cache, connect stan.Conn, subject string) *natsService {
	return &natsService{
		orderRepository: orderRepository,
		cache:           cache,
		connect:         connect,
		subject:         subject,
	}
}

// Subscribe subscribes to a NATS subject and processes incoming messages.
func (ns *natsService) Subscribe(ctx context.Context) error {
	sub, err := ns.connect.Subscribe(ns.subject, ns.process)
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
		return
	}

	id, err := ns.orderRepository.Create(context.Background(), &order)
	if err != nil {
		return
	}

	ns.cache.Set(id, &order)
}

// Publish publishes a message to a NATS subject.
func (ns *natsService) Publish(data []byte) error {
	err := ns.connect.Publish(ns.subject, data)
	if err != nil {
		return fmt.Errorf("can't publish message: %w", err)
	}

	return nil
}
