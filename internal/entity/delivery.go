package entity

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type DeliveryDB struct {
	DeliveryUID string `db:"delivery_uid"`
	Name        string `db:"name"`
	Phone       string `db:"phone"`
	Zip         string `db:"zip"`
	City        string `db:"city"`
	Address     string `db:"address"`
	Region      string `db:"region"`
	Email       string `db:"email"`
}
