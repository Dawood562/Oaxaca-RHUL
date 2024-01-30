package structs

type MenuItem struct {
	MenuItemId string  `json:"menu_item_id" gorm:"column:menuitemid"`
	ItemName   string  `json:"item_name" gorm:"column:itemname"`
	Price      float64 `json:"price" gorm:"column:price"`
	Calories   float64 `json:"calories" gorm:"column:calories"`
}