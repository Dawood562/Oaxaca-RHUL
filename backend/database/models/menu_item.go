package models

type MenuItem struct {
	ID          uint    `json:"itemId" gorm:"primarykey;autoincrement"`
	Name        string  `json:"itemName" gorm:"unique"`
	Description string  `json:"itemDescription"`
	Price       float64 `json:"price"`
	Calories    float64 `json:"calories"`
}
