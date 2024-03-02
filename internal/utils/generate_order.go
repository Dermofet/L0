package utils

import (
	"crypto/rand"
	"log"
	"math/big"
	"time"
)

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SMID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func GenerateOrder() Order {
	trackNumber := generateRandomString(10)
	return Order{
		OrderUID:          generateRandomString(10),
		TrackNumber:       trackNumber,
		Entry:             generateRandomString(4),
		Delivery:          generateRandomDelivery(),
		Payment:           generateRandomPayment(),
		Items:             generateRandomItems(trackNumber),
		Locale:            "en",
		InternalSignature: generateRandomString(10),
		CustomerID:        generateRandomString(5),
		DeliveryService:   generateRandomString(6),
		ShardKey:          generateRandomString(2),
		SMID:              generateRandomNumber(100),
		DateCreated:       time.Now().Format(time.RFC3339), // Use the current time
		OofShard:          generateRandomString(1),
	}
}

// generateRandomString generates a random string of the specified length.
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Fatal(err)
		}
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

// generateRandomNumber generates a random number up to the specified maximum.
func generateRandomNumber(max int) int {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		log.Fatal(err)
	}
	return int(num.Int64())
}

// generateRandomDelivery generates random delivery data.
func generateRandomDelivery() Delivery {
	return Delivery{
		Name:    generateRandomString(10),
		Phone:   generateRandomString(10),
		Zip:     generateRandomString(6),
		City:    generateRandomString(10),
		Address: generateRandomString(20),
		Region:  generateRandomString(10),
		Email:   generateRandomString(10) + "@example.com",
	}
}

// generateRandomPayment generates random payment data.
func generateRandomPayment() Payment {
	return Payment{
		Transaction:  generateRandomString(10),
		RequestID:    generateRandomString(10),
		Currency:     "USD",
		Provider:     generateRandomString(6),
		Amount:       generateRandomNumber(1000),
		PaymentDT:    int(time.Now().Unix()),
		Bank:         generateRandomString(6),
		DeliveryCost: generateRandomNumber(1000),
		GoodsTotal:   generateRandomNumber(1000),
		CustomFee:    generateRandomNumber(100),
	}
}

// generateRandomItems generates random order items.
func generateRandomItems(trackNumber string) []Item {
	var items []Item
	for i := 0; i < 3; i++ {
		items = append(items, Item{
			ChrtID:      generateRandomNumber(100000),
			TrackNumber: trackNumber,
			Price:       generateRandomNumber(1000),
			RID:         generateRandomString(10),
			Name:        generateRandomString(8),
			Sale:        generateRandomNumber(50),
			Size:        generateRandomString(2),
			TotalPrice:  generateRandomNumber(1000),
			NmID:        generateRandomNumber(100000),
			Brand:       generateRandomString(6),
			Status:      generateRandomNumber(300),
		})
	}
	return items
}
