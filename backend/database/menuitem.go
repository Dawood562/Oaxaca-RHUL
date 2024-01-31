package database

type MenuItem struct {
	ID       int     `json:"itemId" gorm:"column:itemid;primary_key"`
	ItemName string  `json:"itemName" gorm:"column:itemname"`
	Price    float64 `json:"price" gorm:"column:price"`
	Calories float64 `json:"calories" gorm:"column:calories"`
}
