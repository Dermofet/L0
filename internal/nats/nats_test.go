package nats

import (
	"context"
	"fmt"
	"testing"
	"time"

	"L0/internal/cache"
	"L0/internal/repository"

	"github.com/golang/mock/gomock"
)

func TestNatsService_Subscribe(t *testing.T) {
	type fields struct {
		orderRepository repository.OrderRepository
		cache           cache.Cache
		connect         *MockConn
		subject         string
		subscription    *MockSubscription
	}
	tests := []struct {
		name    string
		setup   func(f fields)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(f fields) {
				f.connect.EXPECT().Subscribe(f.subject, gomock.Any()).Return(f.subscription, nil)
				f.subscription.EXPECT().Unsubscribe().Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail: can't subscribe",
			setup: func(f fields) {
				f.connect.EXPECT().Subscribe(f.subject, gomock.Any()).Return(nil, fmt.Errorf("subscribe error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				orderRepository: repository.NewMockOrderRepository(ctrl),
				cache:           cache.NewCache(),
				connect:         NewMockConn(ctrl),
				subject:         "test",
				subscription:    NewMockSubscription(ctrl),
			}
			service := NewNatsService(f.orderRepository, f.cache, f.connect, f.subject)
			tt.setup(f)

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := service.Subscribe(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNatsService_Publish(t *testing.T) {
	type fields struct {
		orderRepository repository.OrderRepository
		cache           cache.Cache
		connect         *MockConn
		subject         string
	}
	tests := []struct {
		name    string
		setup   func(f fields)
		wantErr bool
	}{
		{
			name: "success",
			setup: func(f fields) {
				f.connect.EXPECT().Publish(f.subject, []byte("test")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail: can't publish",
			setup: func(f fields) {
				f.connect.EXPECT().Publish(f.subject, []byte("test")).Return(fmt.Errorf("publish error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				orderRepository: repository.NewMockOrderRepository(ctrl),
				cache:           cache.NewCache(),
				connect:         NewMockConn(ctrl),
				subject:         "test",
			}
			service := NewNatsService(f.orderRepository, f.cache, f.connect, f.subject)

			tt.setup(f)

			err := service.Publish([]byte("test"))
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
