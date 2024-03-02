package db

import (
	"L0/internal/entity"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func MustParseTime(layout string, s string) time.Time {
	tt, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}
	return tt
}

// func Test_source_CreateOrder(t *testing.T) {
// 	type fields struct {
// 		db sqlmock.Sqlmock
// 	}
// 	type args struct {
// 		ctx   context.Context
// 		order *entity.Order
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		setup   func(f fields, a args)
// 		wantErr bool
// 	}{
// 		{
// 			name: "CreateOrder_Successful",
// 			args: args{
// 				ctx: context.Background(),
// 				order: &entity.Order{
// 					OrderUID:          "order_uid_1",
// 					TrackNumber:       "track_number_1",
// 					Entry:             "entry_1",
// 					Locale:            "en",
// 					InternalSignature: "",
// 					CustomerID:        "customer_1",
// 					DeliveryService:   "delivery_service_1",
// 					Shardkey:          "shardkey_1",
// 					SmID:              1,
// 					DateCreated:       time.Now(),
// 					OofShard:          "1",
// 					Delivery: entity.Delivery{
// 						Name:    "Test Name",
// 						Phone:   "+1234567890",
// 						Zip:     "12345",
// 						City:    "City",
// 						Address: "Address",
// 						Region:  "Region",
// 						Email:   "test@example.com",
// 					},
// 					Payment: entity.Payment{
// 						Transaction:  "transaction_1",
// 						Currency:     "USD",
// 						Provider:     "provider_1",
// 						Amount:       100,
// 						PaymentDt:    int(time.Now().Unix()),
// 						DeliveryCost: 10,
// 						GoodsTotal:   90,
// 						CustomFee:    0,
// 					},
// 					Items: []entity.Item{
// 						{
// 							ChrtID:      1,
// 							TrackNumber: "track_number_1",
// 							Price:       50,
// 							Rid:         "rid_1",
// 							Name:        "Item 1",
// 							Sale:        0,
// 							Size:        "S",
// 							TotalPrice:  50,
// 							NmID:        1,
// 							Brand:       "Brand 1",
// 							Status:      1,
// 						},
// 					},
// 				},
// 			},
// 			want: "order_uid_1",
// 			setup: func(f fields, a args) {
// 				f.db.ExpectQuery("SELECT * FROM deliveries WHERE phone = $1").
// 					WithArgs(a.order.Delivery.Phone).
// 					WillReturnRows(sqlmock.NewRows([]string{"delivery_uid"}).AddRow("delivery_uid_1"))
// 				f.db.ExpectQuery("SELECT * FROM payments WHERE transaction = $1").
// 					WithArgs(a.order.Payment.Transaction).
// 					WillReturnRows(sqlmock.NewRows([]string{"transaction"}).AddRow("transaction_1"))
// 				f.db.ExpectQuery("SELECT * FROM items WHERE track_number = $1").
// 					WithArgs(a.order.Items[0].TrackNumber).
// 					WillReturnRows(sqlmock.NewRows([]string{"track_number"}).AddRow("track_number_1"))
// 				f.db.ExpectQuery(`INSERT INTO orders
// 					(order_uid, track_number, entry, delivery_uid, payment_transaction, locale, 
// 					internal_signature, customer_id, delivery_service, shardkey, sm_id, 
// 					date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`).
// 					WillReturnRows(sqlmock.NewRows([]string{"order_uid"}).AddRow("order_uid_1"))
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "CreateOrder_DeliveryByPhoneError",
// 			args: args{
// 				ctx: context.Background(),
// 				order: &entity.Order{
// 					Delivery: entity.Delivery{
// 						Phone: "+1234567890",
// 					},
// 				},
// 			},
// 			want:    "",
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
// 			if err != nil {
// 				t.Errorf("can't connect to database: %v", err)
// 				return
// 			}

// 			f := fields{
// 				db: mock,
// 			}

// 			s := &source{
// 				db: sqlx.NewDb(db, "sqlmock"),
// 			}

// 			tt.setup(f, tt.args)

// 			got, err := s.CreateOrder(tt.args.ctx, tt.args.order)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("source.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("source.CreateOrder() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_source_GetOrderByUid(t *testing.T) {
	type fields struct {
		db sqlmock.Sqlmock
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
			name: "ok",
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
				orderRows := sqlmock.NewRows([]string{
					"track_number",
					"entry",
					"delivery_uid",
					"payment_transaction",
					"locale",
					"internal_signature",
					"customer_id",
					"delivery_service",
					"shardkey",
					"sm_id",
					"date_created",
					"oof_shard",
				}).AddRow(
					"WBILMTESTTRACK",                       // track_number
					"WBIL",                                 // entry
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"b563feb7b2b84b6test",                  // payment_transaction
					"en",                                   // locale
					"",                                     // internal_signature
					"test",                                 // customer_id
					"meest",                                // delivery_service
					"9",                                    // shardkey
					99,                                     // sm_id
					MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"), // date_created
					"1", // oof_shard
				)

				deliveryRows := sqlmock.NewRows([]string{
					"delivery_uid",
					"name",
					"phone",
					"zip",
					"city",
					"address",
					"region",
					"email",
				}).AddRow(
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"Test Testov",                          // name
					"+9720000000",                          // phone
					"2639809",                              // zip
					"Kiryat Mozkin",                        // city
					"Ploshad Mira 15",                      // address
					"Kraiot",                               // region
					"test@gmail.com",                       // email
				)

				paymentRows := sqlmock.NewRows([]string{
					"transaction",
					"request_id",
					"currency",
					"provider",
					"amount",
					"payment_dt",
					"bank",
					"delivery_cost",
					"goods_total",
					"custom_fee",
				}).AddRow(
					"b563feb7b2b84b6test", // transaction
					"",                    // request_id
					"USD",                 // currency
					"wbpay",               // provider
					1817,                  // amount
					1637907727,            // payment_dt
					"alpha",               // bank
					1500,                  // delivery_cost
					317,                   // goods_total
					0,                     // custom_fee
				)

				itemRows := sqlmock.NewRows([]string{
					"chrt_id",
					"track_number",
					"price",
					"rid",
					"name",
					"sale",
					"size",
					"total_price",
					"nm_id",
					"brand",
					"status",
				}).AddRow(
					9934930,                 // chrt_id
					"WBILMTESTTRACK",        // track_number
					453,                     // price
					"ab4219087a764ae0btest", // rid
					"Mascaras",              // name
					30,                      // sale
					"0",                     // size
					317,                     // total_price
					2389212,                 // nm_id
					"Vivienne Sabo",         // brand
					202,                     // status
				)

				f.db.ExpectQuery(
					`SELECT
						track_number, entry, delivery_uid, payment_transaction, locale,
						internal_signature, customer_id, delivery_service, shardkey, sm_id,
						date_created, oof_shard
					FROM orders
					WHERE order_uid = $1`,
				).WithArgs(a.uid).WillReturnRows(orderRows)

				f.db.ExpectQuery(
					"SELECT * FROM delivery WHERE delivery_uid = $1",
				).WithArgs("4a6e104d-9d7f-45ff-8de6-37993d709522").WillReturnRows(deliveryRows)

				f.db.ExpectQuery(
					"SELECT * FROM payments WHERE transaction = $1",
				).WithArgs("b563feb7b2b84b6test").WillReturnRows(paymentRows)

				f.db.ExpectQuery(
					"SELECT * FROM items WHERE track_number = $1",
				).WithArgs("WBILMTESTTRACK").WillReturnRows(itemRows)
			},
			wantErr: false,
		},
		{
			name: "fail: can't exec query",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery(
					`BUM BUM`,
				).WithArgs(a.uid).WillReturnError(fmt.Errorf("can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "fail: can't scan rows",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.db.ExpectQuery(
					`SELECT
						track_number, entry, delivery_uid, payment_transaction, locale,
						internal_signature, customer_id, delivery_service, shardkey, sm_id,
						date_created, oof_shard
					FROM orders
					WHERE order_uid = $1`,
				).WithArgs(a.uid).WillReturnRows(sqlmock.NewRows([]string{"track_number"}))
			},
			wantErr: true,
		},
		{
			name: "fail: can't get delivery",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				orderRows := sqlmock.NewRows([]string{
					"track_number",
					"entry",
					"delivery_uid",
					"payment_transaction",
					"locale",
					"internal_signature",
					"customer_id",
					"delivery_service",
					"shardkey",
					"sm_id",
					"date_created",
					"oof_shard",
				}).AddRow(
					"WBILMTESTTRACK",                       // track_number
					"WBIL",                                 // entry
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"b563feb7b2b84b6test",                  // payment_transaction
					"en",                                   // locale
					"",                                     // internal_signature
					"test",                                 // customer_id
					"meest",                                // delivery_service
					"9",                                    // shardkey
					99,                                     // sm_id
					MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"), // date_created
					"1", // oof_shard
				)
				f.db.ExpectQuery(
					`SELECT
						track_number, entry, delivery_uid, payment_transaction, locale,
						internal_signature, customer_id, delivery_service, shardkey, sm_id,
						date_created, oof_shard
					FROM orders
					WHERE order_uid = $1`,
				).WithArgs(a.uid).WillReturnRows(orderRows)
				f.db.ExpectQuery("Bum bum").WillReturnError(fmt.Errorf("can't get delivery: can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "fail: can't get payment",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				orderRows := sqlmock.NewRows([]string{
					"track_number",
					"entry",
					"delivery_uid",
					"payment_transaction",
					"locale",
					"internal_signature",
					"customer_id",
					"delivery_service",
					"shardkey",
					"sm_id",
					"date_created",
					"oof_shard",
				}).AddRow(
					"WBILMTESTTRACK",                       // track_number
					"WBIL",                                 // entry
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"b563feb7b2b84b6test",                  // payment_transaction
					"en",                                   // locale
					"",                                     // internal_signature
					"test",                                 // customer_id
					"meest",                                // delivery_service
					"9",                                    // shardkey
					99,                                     // sm_id
					MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"), // date_created
					"1", // oof_shard
				)

				deliveryRows := sqlmock.NewRows([]string{
					"delivery_uid",
					"name",
					"phone",
					"zip",
					"city",
					"address",
					"region",
					"email",
				}).AddRow(
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"Test Testov",                          // name
					"+9720000000",                          // phone
					"2639809",                              // zip
					"Kiryat Mozkin",                        // city
					"Ploshad Mira 15",                      // address
					"Kraiot",                               // region
					"test@gmail.com",                       // email
				)

				f.db.ExpectQuery(
					`SELECT
						track_number, entry, delivery_uid, payment_transaction, locale,
						internal_signature, customer_id, delivery_service, shardkey, sm_id,
						date_created, oof_shard
					FROM orders
					WHERE order_uid = $1`,
				).WithArgs(a.uid).WillReturnRows(orderRows)

				f.db.ExpectQuery(
					"SELECT * FROM delivery WHERE delivery_uid = $1",
				).WithArgs("4a6e104d-9d7f-45ff-8de6-37993d709522").WillReturnRows(deliveryRows)

				f.db.ExpectQuery("Bum bum").WillReturnError(fmt.Errorf("can't get payment: can't exec query"))
			},
			wantErr: true,
		},
		{
			name: "fail: can't get items",
			args: args{
				ctx: context.Background(),
				uid: "b563feb7b2b84b6test",
			},
			want: nil,
			setup: func(a args, f fields) {
				orderRows := sqlmock.NewRows([]string{
					"track_number",
					"entry",
					"delivery_uid",
					"payment_transaction",
					"locale",
					"internal_signature",
					"customer_id",
					"delivery_service",
					"shardkey",
					"sm_id",
					"date_created",
					"oof_shard",
				}).AddRow(
					"WBILMTESTTRACK",                       // track_number
					"WBIL",                                 // entry
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"b563feb7b2b84b6test",                  // payment_transaction
					"en",                                   // locale
					"",                                     // internal_signature
					"test",                                 // customer_id
					"meest",                                // delivery_service
					"9",                                    // shardkey
					99,                                     // sm_id
					MustParseTime(time.RFC3339, "2021-11-26T06:22:19Z"), // date_created
					"1", // oof_shard
				)

				deliveryRows := sqlmock.NewRows([]string{
					"delivery_uid",
					"name",
					"phone",
					"zip",
					"city",
					"address",
					"region",
					"email",
				}).AddRow(
					"4a6e104d-9d7f-45ff-8de6-37993d709522", // delivery_uid
					"Test Testov",                          // name
					"+9720000000",                          // phone
					"2639809",                              // zip
					"Kiryat Mozkin",                        // city
					"Ploshad Mira 15",                      // address
					"Kraiot",                               // region
					"test@gmail.com",                       // email
				)

				paymentRows := sqlmock.NewRows([]string{
					"transaction",
					"request_id",
					"currency",
					"provider",
					"amount",
					"payment_dt",
					"bank",
					"delivery_cost",
					"goods_total",
					"custom_fee",
				}).AddRow(
					"b563feb7b2b84b6test", // transaction
					"",                    // request_id
					"USD",                 // currency
					"wbpay",               // provider
					1817,                  // amount
					1637907727,            // payment_dt
					"alpha",               // bank
					1500,                  // delivery_cost
					317,                   // goods_total
					0,                     // custom_fee
				)

				f.db.ExpectQuery(
					`SELECT
						track_number, entry, delivery_uid, payment_transaction, locale,
						internal_signature, customer_id, delivery_service, shardkey, sm_id,
						date_created, oof_shard
					FROM orders
					WHERE order_uid = $1`,
				).WithArgs(a.uid).WillReturnRows(orderRows)

				f.db.ExpectQuery(
					"SELECT * FROM delivery WHERE delivery_uid = $1",
				).WithArgs("4a6e104d-9d7f-45ff-8de6-37993d709522").WillReturnRows(deliveryRows)

				f.db.ExpectQuery(
					"SELECT * FROM payments WHERE transaction = $1",
				).WithArgs("b563feb7b2b84b6test").WillReturnRows(paymentRows)

				f.db.ExpectQuery("Bum bum").WillReturnError(fmt.Errorf("can't get items: can't exec query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("can't connect to database: %v", err)
				return
			}
			f := fields{
				db: mock,
			}

			s := &source{
				db: sqlx.NewDb(db, "sqlmock"),
			}

			tt.setup(tt.args, f)
			got, err := s.GetOrderByUid(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("source.GetOrderByUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("source.GetOrderByUid() = %v, want %v", got, tt.want)
			}
		})
	}
}
