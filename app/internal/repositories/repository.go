package repository

import (
	"context"
	"fmt"
	model "postgres-crud/app/internal/models"

	"gorm.io/gorm"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	GetAll(ctx context.Context, page, perPage int) (*[]model.Order, error)
	GetByID(ctx context.Context, id int) (*model.Order, error)
	Update(ctx context.Context, order *model.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		return nil
	})
}

func (r *orderRepository) GetAll(ctx context.Context, page, perPage int) (*[]model.Order, error) {
	var orders []model.Order
	offset := (page - 1) * perPage
	if err := r.db.WithContext(ctx).Preload("OrderItems").Offset(offset).Limit(perPage).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

func (r *orderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order

	if err := r.db.WithContext(ctx).Preload("OrderItems").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Debug().
			Model(&model.Order{}).
			Where("id = ?", order.ID).
			Updates(map[string]interface{}{
				"status": order.Status,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}
