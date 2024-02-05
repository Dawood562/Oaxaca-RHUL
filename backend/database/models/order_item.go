package models

type OrderItem struct {
	OrderID uint     `json:"orderId" gorm:"column:order_id;primarykey"`
	Item    MenuItem `json:"item" gorm:"foreignkey:item_id;primarykey"`
	Notes   string   `json:"notes"`
}
