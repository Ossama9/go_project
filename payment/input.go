package payment

type InputPayment struct {
	ProductId int64   `json:"product_id" building:"required"`
	PricePaid float64 `json:"price_paid" building:"required"`
}
