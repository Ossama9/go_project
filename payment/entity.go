package payment

import "time"

type Payment struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	PricePaid float64   `json:"price_paid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
