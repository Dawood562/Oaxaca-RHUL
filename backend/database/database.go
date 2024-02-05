package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dbUsername, dbName, dbPassword := fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@db:5432/" + dbName
	conn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	db = conn
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	err = db.AutoMigrate(&MenuItem{})
	if err != nil {
		log.Fatal(err)
	}
}

// AddItem adds the given item to the database.
// Returns an error if there is a problem adding the item.
// Item names must be unique.
func AddItem(item *MenuItem) error {
	result := db.Table("menuitem").Create(item)
	return result.Error
}

// EditItem edits the given item with new information
func EditItem(item *MenuItem) error {
	// Check that the item exists
	var count int64
	result := db.Table("menuitem").Where("menuitemid = ?", item.ID).Count(&count)
	if count == 0 {
		return errors.New("Item does not exist")
	}
	// Update the item
	result = db.Table("menuitem").Save(&item)
	return result.Error
}

// RemoveItem removes an item from the menu with the given id
// Returns an error if the item could not be removed
func RemoveItem(id int) error {
	result := db.Table("menuitem").Where("menuitemid = ?", id).Delete(&MenuItem{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("Item with id '%d' does not exist", id))
	}
	return nil
}

// QueryMenu returns the menu items from the database as a slice
// If filter is provided, the returned item slice will be filtered as such
func QueryMenu(filter *MenuFilter) []MenuItem {
	preparedFilter := prepareArgs(filter)

	var data []MenuItem
	db.Table("menuitem").Model(&MenuItem{}).Where("menuItemName LIKE ?", preparedFilter.SearchTerm).Where("calories <= ?", preparedFilter.MaxCalories).Where("price <= ?", preparedFilter.MaxPrice).Find(&data)
	return data
}

// prepareArgs applies defaults to a MenuFilter struct in preparation for use in a query
func prepareArgs(filter *MenuFilter) *MenuFilter {
	ret := &MenuFilter{}

	// If the search term is shorter than 3 chars, disregard it
	if len(filter.SearchTerm) < 3 {
		ret.SearchTerm = "%"
	} else {
		ret.SearchTerm = "%" + filter.SearchTerm + "%"
	}

	// If no max calories are provided, set it to a high number
	if filter.MaxCalories <= 0 {
		ret.MaxCalories = 9999
	} else {
		ret.MaxCalories = filter.MaxCalories
	}

	// If no max price is provided, set it to a high number
	if filter.MaxPrice <= 0 {
		ret.MaxPrice = 9999
	} else {
		ret.MaxPrice = filter.MaxPrice
	}

	return ret
}

/*
Fetches database login details from .env
db_details.txt should be in database folder with following content structure:
<username>
<database name>
<password>
*/
func fetchDBAuth() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		// Load test env if production env not found
		godotenv.Load(".env.test")
	}
	username := os.Getenv("USERNAME")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	return username, dbname, password
}

// ClearMenu clears all items from the menu. For testing use only
func ClearMenu() {
	db.Exec("DELETE FROM menuitem")
}
