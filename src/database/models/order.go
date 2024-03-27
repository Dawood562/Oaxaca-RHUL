package models

import "time"

// Order represents one order in the restaurant system.
type Order struct {
	ID          uint        `json:"orderId" gorm:"primarykey;autoincrement"`
	Time        time.Time   `json:"orderTime" gorm:"default:CURRENT_TIMESTAMP"`
	TableNumber uint        `json:"tableNumber"`
	Bill        float64     `json:"bill"`
	Status      string      `json:"status" gorm:"default:Awaiting Confirmation"`
	Paid        bool        `json:"paid" gorm:"default:false"`
	Items       []OrderItem `json:"items" gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
