package product

type InputProduct struct {
	Name  string  `json:"name" building:"required"`
	Price float64 `json:"price" building:"required"`
}
