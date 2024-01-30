package database

type MenuFilter struct {
	SearchTerm  string  `json:"searchTerm"`
	MaxCalories int     `json:"maxCalories"`
	MaxPrice    float32 `json:"maxPrice"`
}
