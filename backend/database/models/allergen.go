package models

type Allergen struct {
	ID     uint   `gorm:"primarykey" json:"-"`
	ItemID uint   `json:"-"`
	Name   string `json:"name"`
}
