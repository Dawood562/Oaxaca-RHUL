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

var (
	ErrOrderNotFound         error = errors.New("order not found")
	ErrOrderAlreadyPaid      error = errors.New("order already paid for")
	ErrOrderAlreadyConfirmed error = errors.New("order already confirmed")
	ErrOrderAlreadyCancelled error = errors.New("order already cancelled")
)

type OrderStatus string

const (
	StatusAwaitingConfirmation OrderStatus = "Awaiting Confirmation"
	StatusPreparing            OrderStatus = "Preparing"
	StatusDelivered            OrderStatus = "Delivered"
	StatusCancelled            OrderStatus = "Cancelled"
)

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

// FetchItem retrieves the given item from the database
func FetchItem(id int) (models.MenuItem, error) {
	ret := models.MenuItem{}
	res := db.Model(&models.MenuItem{}).Where("ID = ?", id).First(&ret)
	if res.Error != nil {
		return models.MenuItem{}, res.Error
	}
	return ret, nil
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
	// Check that there are no open orders with that table number already
	var count int64
	db.Model(&models.Order{}).Where("status != ? AND table_number = ?", "Complete", item.TableNumber).Count(&count)
	if count > 0 {
		return fmt.Errorf("there is already an open order for table %d", item.TableNumber)
	}
	// Do not re-create menu items when adding them to an OrderItem
	feedback := db.Omit("Items.Item.*").Create(item)
	return feedback.Error
}

func RemoveOrder(id uint) error {
	feedback := db.Delete(&models.Order{ID: id})
	if feedback.Error != nil {
		return feedback.Error
	}
	if feedback.RowsAffected == 0 {
		return errors.New("no items removed from table")
	}
	return nil
}

func FetchOrders(filter ...models.Order) ([]*models.Order, error) {
	dbCopy := *db
	dbLocal := &dbCopy
	var orderData []*models.Order

	if len(filter) > 0 {

		if filter[0].TableNumber > 0 {
			dbLocal = dbLocal.Where("Table_Number = ?", filter[0].TableNumber)
		}
		if len(filter[0].Status) > 0 {
			dbLocal = dbLocal.Where("Status = ?", filter[0].Status)
		}
	}
	dbLocal.Model(&models.Order{}).Preload("Items").Preload("Items.Item").Find(&orderData)

	return orderData, nil
}

// OrderPaid returns true if the order with the given ID has been paid for, returns error if that order does not exist
func OrderPaid(id uint) (bool, error) {
	order := &models.Order{ID: id}
	result := db.Model(order).First(&order)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, ErrOrderNotFound
	}

	return order.Paid, nil
}

// PayOrder updates the payment status of the given order, returns an error if that order cannot be found or is already paid for
func PayOrder(id uint) error {
	paid, err := OrderPaid(id)
	if err != nil {
		return err
	}
	if paid {
		return ErrOrderAlreadyPaid
	}

	// Set order as paid
	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Paid = true
	db.Save(&order)
	return nil
}

// GetOrderStatus returns the Status field of the given order. Returns an error if the order does not exist.
func GetOrderStatus(id uint) (string, error) {
	order := &models.Order{ID: id}
	result := db.Model(order).First(&order)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", errors.New("order not found")
	}

	return order.Status, nil
}

// ConfirmOrder confirms the given order. Returns an error if the order does not exist or is already confirmed.
func ConfirmOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return err
	}
	if status != string(StatusAwaitingConfirmation) {
		return ErrOrderAlreadyConfirmed
	}

	// Set order as paid
	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = string(StatusPreparing)
	db.Save(&order)
	return nil
}

// CancelOrder cancels the given order. Returns an error if the order does not exist or is already cancelled.
func CancelOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return ErrOrderNotFound
	}
	if status == string(StatusCancelled) {
		return ErrOrderAlreadyCancelled
	}

	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = string(StatusCancelled)
	db.Save(&order)

	return nil
}
