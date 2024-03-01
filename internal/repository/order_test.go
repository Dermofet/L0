package repository

import (
	"context"
	"errors"
	"L0/internal/db"
	"L0/internal/entity"
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
