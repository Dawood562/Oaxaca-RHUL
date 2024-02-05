package models

type MenuItem struct {
	ID          int     `json:"itemId" gorm:"column:item_id;primary_key"`
	Name        string  `json:"itemName" gorm:"column:name;unique"`
	Description string  `json:"itemDescription" gorm:"column:description"`
	Price       float64 `json:"price" gorm:"column:price"`
	Calories    float64 `json:"calories" gorm:"column:calories"`
}
