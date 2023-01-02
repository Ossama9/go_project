package product

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(product Product) (Product, error)
	Update(id int, InputProduct InputProduct) (Product, error)
	Delete(id int) error
	GetAll() ([]Product, error)
	GetById(id int) (Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(product Product) (Product, error) {
	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(id int, InputProduct InputProduct) (Product, error) {
	task, err := r.GetById(id)
	if err != nil {
		return task, err
	}

	task.Name = InputProduct.Name
	task.Price = InputProduct.Price

	err = r.db.Save(&task).Error
	if err != nil {
		return task, err
	}

	return task, nil
}

func (r *repository) Delete(id int) error {
	task := &Product{ID: id}
	tx := r.db.Delete(task)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Task not found")
	}

	return nil
}

func (r *repository) GetAll() ([]Product, error) {
	var tasks []Product
	err := r.db.Find(&tasks).Error
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (r *repository) GetById(id int) (Product, error) {
	var product Product

	err := r.db.Where(&Product{ID: id}).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}
