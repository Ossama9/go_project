package main

import (
	"log"
	"project/adapter"
	"project/broadcast"
	"project/handler"
	"project/payment"
	"project/product"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	b := broadcast.NewBroadcaster(20)

	dbURL := "root:password@tcp(127.0.0.1:3306)/go_project?parseTime=true"
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository)

	db.AutoMigrate(&payment.Payment{})
	db.AutoMigrate(&product.Product{})

	ginAdapter := adapter.NewGinAdapter(b, paymentService, productService)

	r := gin.Default()
	r.GET("/stream", ginAdapter.Stream)

	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.GetById)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	r.POST("/payments", ginAdapter.CreatePayment)
	r.GET("/payments", ginAdapter.GetAllPayment)
	r.GET("/payments/:id", ginAdapter.GetPaymentById)
	r.PUT("/payments/:id", ginAdapter.UpdatePayment)
	r.DELETE("/payments/:id", ginAdapter.DeletePayment)

	r.Run()
}
