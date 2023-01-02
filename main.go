package main

import (
	"log"
	"project/handler"
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
	r := gin.Default()
	r.POST("/products", productHandler.Create)
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.GetById)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	r.Run()
}
