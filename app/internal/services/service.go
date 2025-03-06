package service

import (
	"context"
	"fmt"
	constant "postgres-crud/app/internal/constants"
	model "postgres-crud/app/internal/models"
	repository "postgres-crud/app/internal/repositories"
)

type IOrderService interface {
	Create(ctx context.Context, order *model.Order) error
	GetAll(ctx context.Context, page, perPage int) (*[]model.Order, error)
	GetByID(ctx context.Context, id int) (*model.Order, error)
	Update(ctx context.Context, id int, status string) error
}

type orderService struct {
	repo repository.IOrderRepository
}

func NewOrderService(repo repository.IOrderRepository) *orderService {
	return &orderService{repo: repo}
}

func (s *orderService) Create(ctx context.Context, order *model.Order) error {
	order.TotalAmount = calculateTotalAmount(order.OrderItems)
	order.Status = constant.OrderStatusPending

	if err := s.repo.Create(ctx, order); err != nil {
		return err
	}

	return nil
}

func calculateTotalAmount(items []model.OrderItem) float64 {
	totalAmount := 0.00
	for _, item := range items {
		totalAmount += item.Price * float64(item.Quantity)
	}
	return totalAmount
}

func (s *orderService) GetAll(ctx context.Context, page, perPage int) (*[]model.Order, error) {
	return s.repo.GetAll(ctx, page, perPage)
}

func (s *orderService) GetByID(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) Update(ctx context.Context, id int, newStatus string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	allowNewStatusList := getAllowedStatusTransitions(order.Status)

	if !contains(allowNewStatusList, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, newStatus)
	}

	order.Status = newStatus
	return s.repo.Update(ctx, order)
}

func getAllowedStatusTransitions(currentStatus string) []string {
	switch currentStatus {
	case constant.OrderStatusPending:
		return []string{constant.OrderStatusProcessing}
	case constant.OrderStatusProcessing:
		return []string{constant.OrderStatusPaid, constant.OrderStatusCanceled}
	case constant.OrderStatusPaid:
		return []string{constant.OrderStatusShipped}
	case constant.OrderStatusShipped:
		return []string{constant.OrderStatusCompleted, constant.OrderStatusRefunded}
	default:
		return []string{}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
