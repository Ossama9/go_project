package adapter

import (
	"fmt"
	"io"
	"net/http"
	"project/broadcast"
	"project/payment"
	"project/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GinAdapter interface {
	CreatePayment(c *gin.Context)
	UpdatePayment(c *gin.Context)
	DeletePayment(c *gin.Context)
	GetPaymentById(c *gin.Context)
	GetAllPayment(c *gin.Context)

	Stream(c *gin.Context)
}

type Message struct {
	Event string      `json:"data"`
	Data  interface{} `json:"data"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ginAdapter struct {
	broadcast broadcast.Broadcaster
	payment   payment.Service
	product   product.Service
}

func NewGinAdapter(broadcast broadcast.Broadcaster, payment payment.Service, product product.Service) *ginAdapter {

	return &ginAdapter{
		broadcast: broadcast,
		payment:   payment,
		product:   product,
	}
}

func (adapter *ginAdapter) Stream(c *gin.Context) {

	listener := make(chan interface{})

	adapter.broadcast.Register(listener)
	defer adapter.broadcast.Unregister(listener)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				fmt.Println("not a message")
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent(serviceMsg.Event, serviceMsg.Data)
			return true
		}
	})

	fmt.Println("stream is OK")
}

func (adapter *ginAdapter) CreatePayment(c *gin.Context) {
	var input payment.InputPayment
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	item, err := adapter.product.GetById(input.ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Product not found",
			Data:    item,
		})
		return
	}

	newPayment, err := adapter.payment.Create(input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	b := adapter.broadcast

	b.Submit(Message{
		Event: "Payement is created",
		Data:  newPayment,
	})

	response := &Response{
		Success: true,
		Message: "New payment created",
		Data:    newPayment,
	}
	c.JSON(http.StatusCreated, response)

}

func (adapter *ginAdapter) UpdatePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	var input payment.InputPayment
	err = c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	item, err := adapter.product.GetById(input.ProductId)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Product not found",
			Data:    item,
		})
		return
	}

	uPayment, err := adapter.payment.Update(id, input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &Response{
		Success: true,
		Message: "Payement is updated",
		Data:    uPayment,
	}
	c.JSON(http.StatusCreated, response)

	b := adapter.broadcast

	b.Submit(Message{
		Event: "Payement is updated",
		Data:  uPayment,
	})
}

func (adapter *ginAdapter) DeletePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	err = adapter.payment.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Message: "Payment successfully deleted",
	})
}

func (adapter *ginAdapter) GetPaymentById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	payment, err := adapter.payment.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payment,
	})

}

func (adapter *ginAdapter) GetAllPayment(c *gin.Context) {
	payments, err := adapter.payment.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payments,
	})
}
