//go:build integration

package database

import (
	"log"
	"teamproject/database/models"
)

// ClearMenu clears all items from the menu. For testing use only
func ClearMenu() {
	result := db.Where("1 = 1").Delete(&models.MenuItem{})
	if result.Error != nil {
		log.Fatal("Failed to clear menu")
	}
}

// ClearOrders clears all orders from database. For testing use only
func ClearOrders() {
	result := db.Where("1=1").Delete(&models.Order{})
	if result.Error != nil {
		log.Fatal("Failed to clear orders")
	}
}

// CountOrders counts the number of orders currently in the database
func CountOrders() int {
	var count int64
	err := db.Model(&models.Order{}).Where("1=1").Count(&count)
	if err.Error != nil {
		return -1
	}
	return int(count)
}

// FetchOrderItemIDs fetches the IDs of all items in the order with the given Table Number
func FetchOrderItemIDs(t int) []int {
	item := models.Order{}
	db.Model(&models.Order{}).Where("table_number = ?", t).Preload("Items").Preload("Items.Item").First(&item)
	ret := make([]int, len(item.Items))
	for i, item := range item.Items {
		ret[i] = int(item.Item.ID)
	}
	return ret
}

// ResetTestMenu clears the menu table in the database and inserts all testing data
func ResetTestMenu() {
	ClearMenu()

	AddItem(&models.MenuItem{
		ID:          1,
		Name:        "TESTFOOD",
		Description: "Description for TESTFOOD",
		Price:       5.00,
		Calories:    400,
	})
	AddItem(&models.MenuItem{
		ID:          2,
		Name:        "TESTFOOD2",
		Description: "Description for TESTFOOD2",
		Price:       6.00,
		Calories:    500,
	})
	AddItem(&models.MenuItem{
		ID:          3,
		Name:        "TESTFOOD3",
		Description: "Description for TESTFOOD3",
		Price:       7.00,
		Calories:    600,
	})
	AddItem(&models.MenuItem{
		ID:          4,
		Name:        "TESTFOOD4",
		Description: "Description for TESTFOOD4",
		Price:       8.01,
		Calories:    720,
	})
}
