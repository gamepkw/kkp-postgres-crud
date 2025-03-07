package model

type CreateOrder struct {
	CustomerName string             `json:"customerName" validate:"required,max=100"`
	OrderItems   []CreateOrderItems `json:"items" validate:"required,min=1,dive"`
}

type CreateOrderItems struct {
	ProductName string  `json:"productName" validate:"required,max=100"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
	Price       float64 `json:"price" validate:"required,min=0"`
}

func (r *CreateOrder) ToOrder() *Order {
	orderItems := make([]OrderItem, len(r.OrderItems))
	for i, item := range r.OrderItems {
		orderItems[i].ProductName = item.ProductName
		orderItems[i].Quantity = item.Quantity
		orderItems[i].Price = item.Price
	}
	return &Order{
		CustomerName: r.CustomerName,
		OrderItems:   orderItems,
	}
}

type UpdateOrderStatus struct {
	NewStatus string `json:"status" validate:"required,max=20"`
}
