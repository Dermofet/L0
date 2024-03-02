// Package nats provides functionality for interacting with NATS messaging system.

package nats

import "context"

//go:generate mockgen -source=interfaces.go -destination=nats_mock.go -package=nats

//go:generate mockgen -source=C:/Users/pi-77.BARMOUNT/go/pkg/mod/github.com/nats-io/stan.go@v0.10.4/stan.go -destination=stan_mock.go -package=nats
//go:generate mockgen -source=C:/Users/pi-77.BARMOUNT/go/pkg/mod/github.com/nats-io/stan.go@v0.10.4/sub.go -destination=sub_mock.go -package=nats

// NATSService represents a service for handling NATS messaging.
type NATSService interface {
	// Subscribe subscribes to a NATS subject and processes incoming messages.
	// It takes a context, NATS cluster ID, NATS client ID, NATS subject,
	// and NATS URL and returns an error.
	Subscribe(ctx context.Context) error

	// Publish publishes a message to a NATS subject.
	// It takes a NATS cluster ID, NATS client ID, NATS subject,
	// and NATS URL and returns an error.
	Publish(data []byte) error
}
