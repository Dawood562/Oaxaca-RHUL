package database

type MenuItem struct {
	ItemName string  `json:"itemName" gorm:"column:itemname"`
	Price    float64 `json:"price" gorm:"column:price"`
	Calories float64 `json:"calories" gorm:"column:calories"`
}
