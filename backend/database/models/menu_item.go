package models

type MenuItem struct {
	ID          uint    `json:"itemId" gorm:"column:item_id;primarykey"`
	Name        string  `json:"itemName" gorm:"column:name;unique"`
	Description string  `json:"itemDescription" gorm:"column:description"`
	Price       float64 `json:"price" gorm:"column:price"`
	Calories    float64 `json:"calories" gorm:"column:calories"`
}
