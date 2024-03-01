package entity

import (
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type OrderDB struct {
	OrderUID           string    `db:"order_uid"`
	TrackNumber        string    `db:"track_number"`
	Entry              string    `db:"entry"`
	DeliveryUID        string    `db:"delivery_uid"`
	PaymentTransaction string    `db:"payment_transaction"`
	Locale             string    `db:"locale"`
	InternalSignature  string    `db:"internal_signature"`
	CustomerID         string    `db:"customer_id"`
	DeliveryService    string    `db:"delivery_service"`
	Shardkey           string    `db:"shardkey"`
	SmID               int       `db:"sm_id"`
	DateCreated        time.Time `db:"date_created"`
	OofShard           string    `db:"oof_shard"`
}
