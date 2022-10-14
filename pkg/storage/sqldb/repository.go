package sqldb

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(dsn string) (*Storage, error) {
	// dsn := "falfal:Pasword!2@tcp(mysql-dev-db.airy.my.id:3306)/asg_2?charset=utf8mb4&parseTime=True&loc=Local"
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

func (s *Storage) CreateOrder(order *Order) error {
	return s.db.Create(order).Error
}

func (s *Storage) CreateItem(item *Item) error {
	return s.db.Create(item).Error
}

func (s *Storage) GetOrder(orderID uint) (*Order, error) {
	var order Order
	result := s.db.First(&order, orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (s *Storage) GetItem(itemID uint) (*Item, error) {
	var item Item
	result := s.db.First(&item, itemID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (s *Storage) GetItems(orderID uint) ([]Item, error) {
	var items []Item
	result := s.db.Find(&items, "order_id = ?", orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (s *Storage) UpdateOrder(order *Order) error {
	return s.db.Save(order).Error
}

func (s *Storage) UpdateItem(item *Item) error {
	return s.db.Save(item).Error
}

func (s *Storage) DeleteOrder(orderID uint) error {
	return s.db.Delete(&Order{}, orderID).Error
}

func (s *Storage) DeleteItem(itemID uint) error {
	return s.db.Delete(&Item{}, itemID).Error
}
