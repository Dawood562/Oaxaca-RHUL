// Package models contains the data models for each type stored in the database.
package models

// MenuItem represents one item in the restaurant menu.
type MenuItem struct {
	ID          uint       `json:"itemId" gorm:"primarykey;autoincrement"`
	Name        string     `json:"itemName" gorm:"unique"`
	ImageURL    string     `json:"imageURL"`
	Description string     `json:"itemDescription"`
	Price       float64    `json:"price"`
	Calories    float64    `json:"calories"`
	Allergens   []Allergen `json:"allergens" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
