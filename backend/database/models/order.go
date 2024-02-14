package models

import "time"

type Order struct {
	ID          uint        `json:"orderId" gorm:"primarykey"`
	Time        time.Time   `json:"orderTime"`
	TableNumber uint        `json:"tableNumber"`
	Bill        float64     `json:"bill"`
	Status      string      `json:"status"`
	Items       []OrderItem `json:"items" gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
