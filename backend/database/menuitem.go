package database

type MenuItem struct {
	ID          int     `json:"itemId" gorm:"column:menuitemid;primary_key"`
	ItemName    string  `json:"itemName" gorm:"column:menuitemname"`
	Description string  `json:"itemDescription" gorm:"column:itemdescription"`
	Price       float64 `json:"price" gorm:"column:price"`
	Calories    float64 `json:"calories" gorm:"column:calories"`
}
