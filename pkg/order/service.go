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

type Service interface {
	CreateOrder(order Order) (uint, error)
	UpdateOrder(order Order) error
	GetOrder(orderId uint) (*Order, error)
	GetOrders() ([]Order, error)
	DeleteOrder(orderId uint) error
}

type orderService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &orderService{
		repo: repo,
	}
}

func (s *orderService) CreateOrder(order Order) (uint, error) {
	return s.repo.CreateOrder(order)
}

func (s *orderService) GetOrder(orderId uint) (*Order, error) {
	isExist, err := s.repo.OrderExist(orderId)
	if err != nil {
		return nil, err
	}

	if !isExist {
		return nil, errors.New("order not found")
	}
	return s.repo.GetOrder(orderId)
}

func (s *orderService) UpdateOrder(order Order) error {
	isExist, err := s.repo.OrderExist(order.OrderID)
	if err != nil {
		return err
	}

	if !isExist {
		return errors.New("order not found")
	}
	return s.repo.UpdateOrder(order)
}

func (s *orderService) GetOrders() ([]Order, error) {
	return s.repo.GetOrders()
}

func (s *orderService) DeleteOrder(orderId uint) error {
	isExist, err := s.repo.OrderExist(orderId)
	if err != nil {
		return err
	}

	if !isExist {
		return errors.New("order not found")
	}

	return s.repo.DeleteOrder(orderId)
}
