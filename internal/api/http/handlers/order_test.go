package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"L0/internal/entity"
	"L0/internal/usecase"
	"net/http"
	"net/http/httptest"
	reflect "reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

func MustParseTime(layout string, s string) time.Time {
	tt, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return tt
}

func TestOrderHandlers_GetByIdHandler(t *testing.T) {
	type fields struct {
		orderInteractor *usecase.MockOrderInteractor
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name     string
		args     args
		wantBody *entity.Order
		setup    func(a args, f fields)
		wantCode int
	}{
		{
			name: "success GetById usecase",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			wantBody: &entity.Order{
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
				f.orderInteractor.EXPECT().GetByUid(a.ctx, a.uid).Return(order, nil)
			},
			wantCode: 200,
		},
		{
			name: "fail: can't get order by uid from interactor",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			wantBody: nil,
			setup: func(a args, f fields) {
				f.orderInteractor.EXPECT().GetByUid(a.ctx, a.uid).Return(nil, fmt.Errorf("some error"))
			},
			wantCode: 500,
		},
		{
			name: "fail: there is no such order",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			wantBody: nil,
			setup: func(a args, f fields) {
				f.orderInteractor.EXPECT().GetByUid(a.ctx, a.uid).Return(nil, nil)
			},
			wantCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				orderInteractor: usecase.NewMockOrderInteractor(ctrl),
			}
			u := &orderHandlers{
				interactor: f.orderInteractor,
			}

			gin.SetMode(gin.TestMode)
			r := gin.Default()

			r.GET("/orders/id/:uid", u.GetByIdHandler)
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/orders/id/%s", tt.args.uid)
			req, _ := http.NewRequest("GET", url, nil)

			tt.setup(tt.args, f)

			r.ServeHTTP(w, req)
			fmt.Printf("Code = %d\n", w.Code)
			if w.Code != tt.wantCode {
				t.Errorf("orderHandlers.GetOrderByUid() code = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if w.Code != 200 {
				return
			}

			var got entity.Order
			err := json.Unmarshal(w.Body.Bytes(), &got)
			if err != nil {
				t.Errorf("can't unmarshal body: %v", err)
				return
			}

			if !reflect.DeepEqual(&got, tt.wantBody) {
				t.Errorf("userInteractor.GetUser() = %v, want %v", &got, tt.wantBody)
			}
		})
	}
}
