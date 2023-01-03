package main

import (
	"log"
	"project/handler"
	"project/payment"
	"project/product"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dbURL := "root:password@tcp(127.0.0.1:3306)/go_project"
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository)
	paymenttHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()

	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.GetById)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	r.POST("/payments", paymenttHandler.Create)
	r.GET("/payments", paymenttHandler.GetAll)
	r.GET("/payments/:id", paymenttHandler.GetById)
	r.PUT("/payments/:id", paymenttHandler.Update)
	r.DELETE("/payments/:id", paymenttHandler.Delete)

	r.Run()
}
