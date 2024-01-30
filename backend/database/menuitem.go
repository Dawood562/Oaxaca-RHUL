package database

type MenuItem struct {
	ItemId   uint    `json:"itemId" gorm:"column:itemid;primary_key;auto_increment;not_null"`
	ItemName string  `json:"itemName" gorm:"column:itemname"`
	Price    float64 `json:"price" gorm:"column:price"`
	Calories float64 `json:"calories" gorm:"column:calories"`
}
