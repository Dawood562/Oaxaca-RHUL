//go:build integration

package database

import (
	"log"
)

// ClearMenu clears all items from the menu. For testing use only
func ClearMenu() {
	result := db.Where("1 = 1").Delete(&MenuItem{})
	if result.Error != nil {
		log.Fatal("Failed to clear database after test")
	}
}

// ResetTestMenu clears the menu table in the database and inserts all testing data
func ResetTestMenu() {
	ClearMenu()

	AddItem(&MenuItem{
		ID:          1,
		Name:        "TESTFOOD",
		Description: "Description for TESTFOOD",
		Price:       5.00,
		Calories:    400,
	})
	AddItem(&MenuItem{
		ID:          2,
		Name:        "TESTFOOD2",
		Description: "Description for TESTFOOD2",
		Price:       6.00,
		Calories:    500,
	})
	AddItem(&MenuItem{
		ID:          3,
		Name:        "TESTFOOD3",
		Description: "Description for TESTFOOD3",
		Price:       7.00,
		Calories:    600,
	})
	AddItem(&MenuItem{
		ID:          4,
		Name:        "TESTFOOD4",
		Description: "Description for TESTFOOD4",
		Price:       8.01,
		Calories:    720,
	})
}
