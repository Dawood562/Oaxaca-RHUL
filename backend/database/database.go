package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"teamproject/database/models"

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

	// AutoMigrate all tables
	err = db.AutoMigrate(&models.MenuItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.OrderItem{})
	if err != nil {
		log.Fatal(err)
	}
}

// AddItem adds the given item to the database.
// Returns an error if there is a problem adding the item.
// Item names must be unique.
func AddItem(item *models.MenuItem) error {
	result := db.Create(item)
	return result.Error
}

// EditItem edits the given item with new information
func EditItem(item *models.MenuItem) error {
	// Check that the item exists
	result := db.First(&models.MenuItem{ID: item.ID})
	if result.RowsAffected == 0 {
		return errors.New("item does not exist")
	}
	// Update the item
	result = db.Save(&item)
	return result.Error
}

// RemoveItem removes an item from the menu with the given id
// Returns an error if the item could not be removed
func RemoveItem(id uint) error {
	result := db.Delete(&models.MenuItem{ID: id})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("item with id '%d' does not exist", id)
	}
	return nil
}

// QueryMenu returns the menu items from the database as a slice
// If filter is provided, the returned item slice will be filtered as such
func QueryMenu(filter *MenuFilter) []models.MenuItem {
	preparedFilter := prepareArgs(filter)

	var data []models.MenuItem
	db.Model(&models.MenuItem{}).Where("name LIKE ?", preparedFilter.SearchTerm).Where("calories <= ?", preparedFilter.MaxCalories).Where("price <= ?", preparedFilter.MaxPrice).Find(&data)
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

func AddOrder(item *models.Order) error {
	feedback := db.Create(item)
	return feedback.Error
}

func RemoveOrder(id uint) error {
	err := removeAllOrderChildItems(id)
	if err != nil {
		return err
	}

	feedback := db.Delete(&models.Order{ID: id})
	if feedback.Error != nil {
		return feedback.Error
	}
	if feedback.RowsAffected == 0 {
		return errors.New("no items removed from table")
	}
	return nil
}

func removeAllOrderChildItems(id uint) error {
	allOrderItems := fetchOrderItems()

	for _, items := range allOrderItems {
		if id == items.OrderID {
			err := db.Delete(&models.OrderItem{OrderID: id})
			if err.Error != nil || err.RowsAffected == 0 {
				return errors.New("Did not remove child items?!?!?")
			}
		}
	}
	return nil
}

func FetchOrders(filter ...models.Order) ([]models.Order, error) {
	dbCopy := *db
	dbLocal := &dbCopy
	var orderData []models.Order
	var orderItemData = fetchOrderItems()

	if len(filter) > 0 {

		if filter[0].TableNumber > 0 {
			dbLocal = dbLocal.Where("Table_Number = ?", filter[0].TableNumber)
		}
		if len(filter[0].Status) > 0 {
			dbLocal = dbLocal.Where("Status = ?", filter[0].Status)
		}
	}
	dbLocal.Model(&orderData).Find(&orderData)

	// Iterate through each order and append its order items into Items field
	for i, oData := range orderData {
		for _, oiData := range orderItemData {
			if oData.ID == oiData.OrderID {
				beforeSize := len(orderData[i].Items)
				orderData[i].Items = append(orderData[i].Items, oiData)
				if len(orderData[i].Items) <= beforeSize {
					return nil, errors.New("OrderItem was not appended to the order successfully")
				}
			}
		}
	}

	return orderData, nil
}

func fetchOrderItems() []models.OrderItem {
	var data []models.OrderItem
	db.Model(&data).Find(&data)
	return data
}
