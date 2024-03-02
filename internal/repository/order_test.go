package repository

import (
	"L0/internal/db"
	"L0/internal/entity"
	"context"
	"errors"
	"fmt"
	reflect "reflect"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func MustParseTime(layout string, s string) time.Time {
	tt, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return tt
}

func TestOrderRepository_Create(t *testing.T) {
	type fields struct {
		source *db.MockOrderSource
	}
	type args struct {
		ctx   context.Context
		order *entity.Order
	}
	tests := []struct {
		name    string
		args    args
		setup   func(args, fields)
		wantID  string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				order: &entity.Order{},
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().CreateOrder(a.ctx, a.order).Return("generated_id", nil)
			},
			wantID:  "generated_id",
			wantErr: false,
		},
		{
			name: "fail: can't create order",
			args: args{
				ctx:   context.Background(),
				order: &entity.Order{},
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().CreateOrder(a.ctx, a.order).Return("", errors.New("create error"))
			},
			wantID:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockOrderSource(ctrl),
			}
			repo := &orderRepository{
				source: f.source,
			}

			tt.setup(tt.args, f)

			gotID, err := repo.Create(tt.args.ctx, tt.args.order)

			if (err != nil) != tt.wantErr {
				t.Errorf("orderRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("orderRepository.Create() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func Test_GetByUid(t *testing.T) {
	type fields struct {
		source *db.MockOrderSource
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Order
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: GetById userRepository",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: &entity.Order{
				OrderUID:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: entity.Delivery{
					Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Kiryat Mozkin",
					Address: "Ploshad Mira 15",
					Region:  "Kraiot",
					Email:   "test@gmail.com",
				},
				Payment: entity.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestID:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDt:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []entity.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Price:       453,
						Rid:         "ab4219087a764ae0btest",
						Name:        "Mascaras",
						Sale:        30,
						Size:        "0",
						TotalPrice:  317,
						NmID:        2389212,
						Brand:       "Vivienne Sabo",
						Status:      202,
					},
				},
				Locale:            "en",
				InternalSignature: "",
				CustomerID:        "test",
				DeliveryService:   "meest",
				Shardkey:          "9",
				SmID:              99,
				DateCreated:       MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"),
				OofShard:          "1",
			},
			setup: func(a args, f fields) {
				res := &entity.Order{
					OrderUID:    "b563feb7b2b84b6test",
					TrackNumber: "WBILMTESTTRACK",
					Entry:       "WBIL",
					Delivery: entity.Delivery{
						Name:    "Test Testov",
						Phone:   "+9720000000",
						Zip:     "2639809",
						City:    "Kiryat Mozkin",
						Address: "Ploshad Mira 15",
						Region:  "Kraiot",
						Email:   "test@gmail.com",
					},
					Payment: entity.Payment{
						Transaction:  "b563feb7b2b84b6test",
						RequestID:    "",
						Currency:     "USD",
						Provider:     "wbpay",
						Amount:       1817,
						PaymentDt:    1637907727,
						Bank:         "alpha",
						DeliveryCost: 1500,
						GoodsTotal:   317,
						CustomFee:    0,
					},
					Items: []entity.Item{
						{
							ChrtID:      9934930,
							TrackNumber: "WBILMTESTTRACK",
							Price:       453,
							Rid:         "ab4219087a764ae0btest",
							Name:        "Mascaras",
							Sale:        30,
							Size:        "0",
							TotalPrice:  317,
							NmID:        2389212,
							Brand:       "Vivienne Sabo",
							Status:      202,
						},
					},
					Locale:            "en",
					InternalSignature: "",
					CustomerID:        "test",
					DeliveryService:   "meest",
					Shardkey:          "9",
					SmID:              99,
					DateCreated:       MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"),
					OofShard:          "1",
				}
				f.source.EXPECT().GetOrderByUid(a.ctx, a.uid).Return(res, nil)
			},
			wantErr: false,
		},
		{
			name: "fail: can't get order by uid from db",
			args: args{
				ctx: context.Background(),
				uid: "test",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.source.EXPECT().GetOrderByUid(a.ctx, a.uid).Return(nil, errors.New("test"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockOrderSource(ctrl),
			}

			r := NewOrderRepository(f.source)

			tt.setup(tt.args, f)

			got, err := r.GetByUid(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_GetAll(t *testing.T) {
	type fields struct {
		source *db.MockOrderSource
	}
	tests := []struct {
		name    string
		setup   func(f fields)
		want    []*entity.Order
		wantErr bool
	}{
		{
			name: "success",
			setup: func(f fields) {
				orders := []*entity.Order{
					{OrderUID: "b563feb7b2b84b6test1"},
					{OrderUID: "b563feb7b2b84b6test2"},
				}
				f.source.EXPECT().GetAllOrders(gomock.Any()).Return(orders, nil)
			},
			want:    []*entity.Order{{OrderUID: "b563feb7b2b84b6test1"}, {OrderUID: "b563feb7b2b84b6test2"}},
			wantErr: false,
		},
		{
			name: "fail: can't get orders",
			setup: func(f fields) {
				f.source.EXPECT().GetAllOrders(gomock.Any()).Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockOrderSource(ctrl),
			}
			repo := NewOrderRepository(f.source)
			tt.setup(f)

			got, err := repo.GetAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_Delete(t *testing.T) {
	type fields struct {
		source *db.MockOrderSource
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(args, fields)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				uid: "test_uid",
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().DeleteOrder(a.ctx, a.uid).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "fail: can't delete order",
			args: args{
				ctx: context.Background(),
				uid: "test_uid",
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().DeleteOrder(a.ctx, a.uid).Return(errors.New("delete error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockOrderSource(ctrl),
			}
			repo := &orderRepository{
				source: f.source,
			}

			tt.setup(tt.args, f)

			err := repo.Delete(tt.args.ctx, tt.args.uid)

			if (err != nil) != tt.wantErr {
				t.Errorf("orderRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
