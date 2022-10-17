package sqldb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"fga-asg-2/pkg/order"
)

type Item struct {
	ItemID      uint `gorm:"primary_key;auto_increment;not_null;unique"`
	ItemCode    string
	Description string
	Quantity    uint
	OrderId     uint
}

type Order struct {
	OrderID      uint `gorm:"primary_key;auto_increment;not_null;unique"`
	CustomerName string
	OrderedAt    string
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Item{})
	db.AutoMigrate(&Order{})

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *Storage) CreateOrder(order order.Order) (uint, error) {
	var orderId uint
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		dbOrder := Order{
			CustomerName: order.CustomerName,
			OrderedAt:    order.OrderedAt,
		}
		if err := tx.Create(&dbOrder).Error; err != nil {
			return err
		}

		// Create items
		for _, item := range order.Items {
			dbItem := Item{
				ItemCode:    item.ItemCode,
				Description: item.Description,
				Quantity:    item.Quantity,
				OrderId:     dbOrder.OrderID,
			}
			if err := tx.Create(&dbItem).Error; err != nil {
				return err
			}
		}

		// Order ID of the created order
		orderId = dbOrder.OrderID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (s *Storage) UpdateOrder(order order.Order) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update order
		dbOrder := Order{
			OrderID:      order.OrderID,
			CustomerName: order.CustomerName,
			OrderedAt:    order.OrderedAt,
		}
		if err := tx.Save(&dbOrder).Error; err != nil {
			return err
		}

		// Delete items
		if err := tx.Delete(&Item{}, "order_id = ?", order.OrderID).Error; err != nil {
			return err
		}

		// Create items
		for _, item := range order.Items {
			dbItem := Item{
				ItemCode:    item.ItemCode,
				Description: item.Description,
				Quantity:    item.Quantity,
				OrderId:     dbOrder.OrderID,
			}
			if err := tx.Create(&dbItem).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetOrder(orderID uint) (*order.Order, error) {
	var dbOrder Order

	// Get order
	result := s.db.First(&dbOrder, orderID)
	if result.Error != nil {
		return nil, result.Error
	}

	// Get items
	dbItems, err := s.GetItems(orderID)
	if err != nil {
		return nil, err
	}

	items := make([]order.Item, len(dbItems))
	for i, dbItem := range dbItems {
		items[i] = order.Item{
			ItemCode:    dbItem.ItemCode,
			Description: dbItem.Description,
			Quantity:    dbItem.Quantity,
		}
	}

	return &order.Order{
		OrderID:      dbOrder.OrderID,
		CustomerName: dbOrder.CustomerName,
		OrderedAt:    dbOrder.OrderedAt,
		Items:        items,
	}, nil
}

func (s *Storage) GetOrders() ([]order.Order, error) {
	var dbOrders []Order

	// Get orders
	result := s.db.Find(&dbOrders)
	if result.Error != nil {
		return nil, result.Error
	}

	orders := make([]order.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		dbItems, err := s.GetItems(dbOrder.OrderID)
		if err != nil {
			return nil, err
		}

		// Get items for each order
		items := make([]order.Item, len(dbItems))
		for j, dbItem := range dbItems {
			items[j] = order.Item{
				ItemCode:    dbItem.ItemCode,
				Description: dbItem.Description,
				Quantity:    dbItem.Quantity,
			}
		}

		orders[i] = order.Order{
			OrderID:      dbOrder.OrderID,
			CustomerName: dbOrder.CustomerName,
			OrderedAt:    dbOrder.OrderedAt,
			Items:        items,
		}
	}

	return orders, nil
}

func (s *Storage) GetItems(orderID uint) ([]Item, error) {
	var items []Item
	result := s.db.Find(&items, "order_id = ?", orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (s *Storage) DeleteOrder(orderID uint) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// delete items
		if err := tx.Delete(&Item{}, "order_id = ?", orderID).Error; err != nil {
			return err
		}

		// delete order
		if err := tx.Delete(&Order{}, orderID).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) OrderExist(orderID uint) (bool, error) {
	var count int64
	result := s.db.Model(&Order{}).Where("order_id = ?", orderID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
