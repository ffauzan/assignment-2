package order

import "errors"

type Order struct {
	OrderID      uint
	CustomerName string
	OrderedAt    string
	Items        []Item
}

type Item struct {
	ItemCode    string
	Description string
	Quantity    uint
}

type Repository interface {
	CreateOrder(order Order) (uint, error)
	UpdateOrder(order Order) error
	GetOrder(orderId uint) (*Order, error)
	GetOrders() ([]Order, error)
	DeleteOrder(orderId uint) error
	OrderExist(orderId uint) (bool, error)
}

type OrderService struct {
	repo Repository
}

func NewService(repo Repository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(order Order) (uint, error) {
	return s.repo.CreateOrder(order)
}

func (s *OrderService) GetOrder(orderId uint) (*Order, error) {
	isExist, err := s.repo.OrderExist(orderId)
	if err != nil {
		return nil, err
	}

	if !isExist {
		return nil, errors.New("order not found")
	}
	return s.repo.GetOrder(orderId)
}

func (s *OrderService) UpdateOrder(order Order) error {
	isExist, err := s.repo.OrderExist(order.OrderID)
	if err != nil {
		return err
	}

	if !isExist {
		return errors.New("order not found")
	}
	return s.repo.UpdateOrder(order)
}

func (s *OrderService) GetOrders() ([]Order, error) {
	return s.repo.GetOrders()
}

func (s *OrderService) DeleteOrder(orderId uint) error {
	isExist, err := s.repo.OrderExist(orderId)
	if err != nil {
		return err
	}

	if !isExist {
		return errors.New("order not found")
	}

	return s.repo.DeleteOrder(orderId)
}
