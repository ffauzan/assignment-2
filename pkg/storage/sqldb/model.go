package sqldb

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ItemID      uint
	ItemCode    string
	Description string
	Quantity    uint
	OrderId     uint
}

type Order struct {
	gorm.Model
	OrderID      uint
	CustomerName string
	OrderedAt    time.Time
}
