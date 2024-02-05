package util

import "teamproject/database"

// ResetTestMenu clears the menu table in the database and inserts all testing data
func ResetTestMenu() {
	database.ClearMenu()

	database.AddItem(&database.MenuItem{
		ID:          1,
		Name:        "TESTFOOD",
		Description: "Description for TESTFOOD1",
		Price:       5.00,
		Calories:    400,
	})
	database.AddItem(&database.MenuItem{
		ID:          2,
		Name:        "TESTFOOD2",
		Description: "Description for TESTFOOD2",
		Price:       6.00,
		Calories:    500,
	})
	database.AddItem(&database.MenuItem{
		ID:          3,
		Name:        "TESTFOOD3",
		Description: "Description for TESTFOOD3",
		Price:       7.00,
		Calories:    600,
	})
	database.AddItem(&database.MenuItem{
		ID:          4,
		Name:        "TESTFOOD4",
		Description: "Description for TESTFOOD4",
		Price:       8.01,
		Calories:    720,
	})
}
