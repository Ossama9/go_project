package payment

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(payment Payment) (Payment, error)
	Update(id int, InputPayment InputPayment) (Payment, error)
	Delete(id int) error
	GetAll() ([]Payment, error)
	GetById(id int) (Payment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(payment Payment) (Payment, error) {
	err := r.db.Create(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Update(id int, InputPayment InputPayment) (Payment, error) {
	payment, err := r.GetById(id)
	if err != nil {
		return payment, err
	}

	payment.ProductId = InputPayment.ProductId
	payment.PricePaid = InputPayment.PricePaid

	err = r.db.Save(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Delete(id int) error {
	payment := &Payment{Id: id}
	tx := r.db.Delete(payment)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Task not found")
	}

	return nil
}

func (r *repository) GetAll() ([]Payment, error) {
	var payments []Payment
	err := r.db.Find(&payments).Error
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func (r *repository) GetById(id int) (Payment, error) {
	var payment Payment

	err := r.db.Where(&Payment{Id: id}).First(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}
