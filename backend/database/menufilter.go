package database

type MenuFilter struct {
	LowCalories bool
	MaxPrice    float32
	Allergens   []string
}
