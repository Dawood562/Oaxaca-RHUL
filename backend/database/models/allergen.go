package models

type Allergen struct {
	ID     uint     `gorm:"primarykey" json:"-"`
	Item   MenuItem `json:"-"`
	ItemID uint     `json:"-"`
	Name   string   `json:"name"`
}
