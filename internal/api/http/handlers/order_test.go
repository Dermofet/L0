package handlers

import (
	"L0/internal/entity"
	"L0/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
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
			defer ctrl.Finish()
			f := fields{
				orderInteractor: usecase.NewMockOrderInteractor(ctrl),
			}
			h := &orderHandlers{
				interactor: f.orderInteractor,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "uid", Value: tt.args.uid}}

			tt.setup(tt.args, f)

			h.GetByIdHandler(c)

			if w.Code != tt.wantCode {
				t.Errorf("orderHandlers.GetOrderByUid() code = %v, wantCode %v", w.Code, tt.wantCode)
				return
			}

			if w.Code != http.StatusOK {
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

func TestOrderHandlers_DeleteHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockOrderInteractor
	}
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name     string
		args     args
		setup    func(f fields)
		wantCode int
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			wantCode: http.StatusNoContent,
			setup: func(f fields) {
				f.interactor.EXPECT().Delete(gomock.Any(), "b563feb7b2b84b6test").Return(nil)
			},
		},
		{
			name: "fail: can't delete order",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			wantCode: http.StatusInternalServerError,
			setup: func(f fields) {
				f.interactor.EXPECT().Delete(gomock.Any(), "b563feb7b2b84b6test").Return(fmt.Errorf("some error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				interactor: usecase.NewMockOrderInteractor(ctrl),
			}
			h := &orderHandlers{
				interactor: f.interactor,
			}
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = gin.Params{{Key: "uid", Value: tt.args.uid}}

			tt.setup(f)

			h.DeleteHandler(c)

			if c.Writer.Status() != tt.wantCode {
				t.Errorf("DeleteHandler() code = %v, wantCode %v", c.Writer.Status(), tt.wantCode)
			}
		})
	}
}

func TestOrderHandlers_GetAllHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockOrderInteractor
	}
	tests := []struct {
		name     string
		setup    func(f fields)
		wantCode int
		wantBody interface{}
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			setup: func(f fields) {
				orders := []*entity.Order{
					{OrderUID: "b563feb7b2b84b6test1"},
					{OrderUID: "b563feb7b2b84b6test2"},
				}
				f.interactor.EXPECT().GetAll(gomock.Any()).Return(orders, nil)
			},
			wantBody: []string{"b563feb7b2b84b6test1", "b563feb7b2b84b6test2"},
		},
		{
			name:     "fail: can't get orders",
			wantCode: http.StatusInternalServerError,
			setup: func(f fields) {
				f.interactor.EXPECT().GetAll(gomock.Any()).Return(nil, fmt.Errorf("some error"))
			},
			wantBody: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				interactor: usecase.NewMockOrderInteractor(ctrl),
			}
			h := &orderHandlers{
				interactor: f.interactor,
			}
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tt.setup(f)

			h.GetAllHandler(c)

			if w.Code != tt.wantCode {
				t.Errorf("GetAllHandler() code = %v, wantCode %v", c.Writer.Status(), tt.wantCode)
			}

			if w.Code != http.StatusOK {
				return
			}

			var got []string
			fmt.Printf("body: %s\n", w.Body.String())
			err := json.Unmarshal(w.Body.Bytes(), &got)
			if err != nil {
				t.Errorf("can't unmarshal body: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.wantBody) {
				t.Errorf("userInteractor.GetUser() = %v, want %v", &got, tt.wantBody)
			}
		})
	}
}