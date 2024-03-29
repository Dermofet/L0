package usecase

import (
	"L0/internal/cache"
	"L0/internal/entity"
	"L0/internal/repository"
	"context"
	"fmt"
	"reflect"
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

func Test_GetByUid(t *testing.T) {
	type fields struct {
		orderRepository *repository.MockOrderRepository
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
			name: "success GetById usecase",
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
				order := &entity.Order{
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
				f.orderRepository.EXPECT().GetByUid(a.ctx, a.uid).Return(order, nil)
			},
			wantErr: false,
		},
		{
			name: "fail: can't get order by uid from repository",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.orderRepository.EXPECT().GetByUid(a.ctx, a.uid).Return(nil, fmt.Errorf("can't get order by uid from repository"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				orderRepository: repository.NewMockOrderRepository(ctrl),
			}
			u := &orderInteractor{
				repo: f.orderRepository,
			}

			tt.setup(tt.args, f)

			got, err := u.GetByUid(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderInteractor_GetAll(t *testing.T) {
	type fields struct {
		orderRepository *repository.MockOrderRepository
		cache           *cache.MockCache
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		setup   func(f fields, a args)
		want    []*entity.Order
		wantErr bool
	}{
		{
			name: "success: from cache",
			args: args{
				ctx: context.Background(),
			},
			setup: func(f fields, a args) {
				orders := []*entity.Order{
					{
						OrderUID: "b563feb7b2b84b6test1",
					},
					{
						OrderUID: "b563feb7b2b84b6test2",
					},
				}
				interfaceSlice := make([]interface{}, len(orders))
				for i, v := range orders {
					interfaceSlice[i] = v
				}
				f.cache.EXPECT().GetAll().Return(interfaceSlice, true)
			},
			want:    []*entity.Order{{OrderUID: "b563feb7b2b84b6test1"}, {OrderUID: "b563feb7b2b84b6test2"}},
			wantErr: false,
		},
		{
			name: "success: from db",
			args: args{
				ctx: context.Background(),
			},
			setup: func(f fields, a args) {
				orders := []*entity.Order{
					{
						OrderUID: "b563feb7b2b84b6test1",
					},
					{
						OrderUID: "b563feb7b2b84b6test2",
					},
				}
				f.cache.EXPECT().GetAll().Return([]interface{}{}, false)
				f.orderRepository.EXPECT().GetAll(a.ctx).Return(orders, nil)
			},
			want:    []*entity.Order{{OrderUID: "b563feb7b2b84b6test1"}, {OrderUID: "b563feb7b2b84b6test2"}},
			wantErr: false,
		},
		{
			name: "fail: can't get orders",
			args: args{
				ctx: context.Background(),
			},
			setup: func(f fields, a args) {
				f.cache.EXPECT().GetAll().Return(nil, false)
				f.orderRepository.EXPECT().GetAll(a.ctx).Return(nil, fmt.Errorf("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				orderRepository: repository.NewMockOrderRepository(ctrl),
				cache:           cache.NewMockCache(ctrl),
			}
			u := &orderInteractor{
				repo:  f.orderRepository,
				cache: f.cache,
			}

			tt.setup(f, tt.args)

			got, err := u.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.GetUser() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestOrderInteractor_Delete(t *testing.T) {
	type fields struct {
		orderRepository *repository.MockOrderRepository
		cache           *cache.MockCache
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(f fields, a args)
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			setup: func(f fields, a args) {
				f.orderRepository.EXPECT().Delete(a.ctx, a.uid).Return(nil)
				f.cache.EXPECT().Delete(a.uid)
			},
			wantErr: false,
		},
		{
			name: "fail: can't delete order",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			setup: func(f fields, a args) {
				f.orderRepository.EXPECT().Delete(a.ctx, a.uid).Return(fmt.Errorf("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				orderRepository: repository.NewMockOrderRepository(ctrl),
				cache:           cache.NewMockCache(ctrl),
			}
			u := &orderInteractor{
				repo:  f.orderRepository,
				cache: f.cache,
			}

			tt.setup(f, tt.args)

			err := u.Delete(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
