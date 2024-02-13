package models

import "time"

type Order struct {
	ID          uint        `json:"orderId" gorm:"column:order_id;primarykey"`
	Time        time.Time   `json:"orderTime" gorm:"column:order_time"`
	TableNumber uint        `json:"tableNumber" gorm:"column:table_number"`
	Bill        float64     `json:"bill" gorm:"column:bill"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items" gorm:"foreignkey:order_id"`
}
