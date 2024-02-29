package models

type OrderItem struct {
	ID        uint `gorm:"primarykey"`
	OrderID   uint `json:"orderId"`
	ItemRefer uint
	Item      MenuItem `json:"itemId" gorm:"foreignkey:ItemRefer;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Notes     string   `json:"notes"`
}
