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
	ErrOrderAlreadyReady     error = errors.New("order already ready")
	ErrOrderAlreadyDelivered error = errors.New("order already delivered")
)

const (
	StatusAwaitingConfirmation string = "Awaiting Confirmation"
	StatusPreparing            string = "Preparing"
	StatusReady                string = "Ready"
	StatusDelivered            string = "Delivered"
	StatusCancelled            string = "Cancelled"
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
	err = db.AutoMigrate(&models.Allergen{})
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
	// Delete old allergen relations
	db.Where(&models.Allergen{ItemID: item.ID}).Delete(&models.Allergen{})
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
	var data []models.MenuItem
	allergens := filterAllergens(filter.Allergens)
	dbLocal := db.Model(&models.MenuItem{}).Preload("Allergens")
	if len(filter.SearchTerm) > 3 {
		dbLocal = dbLocal.Where("LOWER(menu_items.name) LIKE LOWER(?)", fmt.Sprintf("%%%s%%", filter.SearchTerm))
	}
	if filter.MaxCalories > 0 {
		dbLocal = dbLocal.Where("calories <= ?", filter.MaxCalories)
	}
	if filter.MaxPrice > 0 {
		dbLocal = dbLocal.Where("price <= ?", filter.MaxPrice)
	}
	if len(allergens) > 0 {
		subQuery := db.Model(&models.Allergen{}).Where("LOWER(allergens.name) IN ?", allergens).Group("allergens.item_id").Select("COUNT(*) as num_allergens, allergens.item_id as item_id")
		dbLocal = dbLocal.Joins("FULL OUTER JOIN (?) as allergen_count ON allergen_count.item_id = menu_items.id", subQuery).Where("allergen_count.num_allergens = 0 OR allergen_count.num_allergens is NULL")
	}

	dbLocal.Find(&data)
	return data
}

func filterAllergens(allergens []string) []string {
	ret := []string{}
	for _, allergen := range allergens {
		if len(allergen) > 0 {
			ret = append(ret, allergen)
		}
	}
	return ret
}

// FetchItem retrieves the given item from the database
func FetchItem(id int) (models.MenuItem, error) {
	ret := models.MenuItem{}
	res := db.Model(&models.MenuItem{}).Preload("Allergens").Where("ID = ?", id).First(&ret)
	if res.Error != nil {
		return models.MenuItem{}, res.Error
	}
	return ret, nil
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

func FetchOrders(confirmed bool) ([]*models.Order, error) {
	dbCopy := *db
	dbLocal := &dbCopy
	var orderData []*models.Order

	if confirmed {
		dbLocal = dbLocal.Where("Status != ?", StatusAwaitingConfirmation)
	}
	dbLocal.Model(&models.Order{}).Preload("Items").Preload("Items.Item").Preload("Items.Item.Allergens").Find(&orderData)

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
	db.Save(order)
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
		return "", ErrOrderNotFound
	}

	return order.Status, nil
}

// ConfirmOrder confirms the given order. Returns an error if the order does not exist or is already confirmed.
func ConfirmOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return err
	}
	if status != StatusAwaitingConfirmation {
		return ErrOrderAlreadyConfirmed
	}

	// Set order as paid
	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = StatusPreparing
	db.Save(order)
	return nil
}

// CancelOrder cancels the given order. Returns an error if the order does not exist or is already cancelled.
func CancelOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return ErrOrderNotFound
	}
	if status == StatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = StatusCancelled
	db.Save(order)

	return nil
}

// ReadyOrder marks the given order as ready for delivery. Returns an error if the order is cancelled, already ready or already delivered.
func ReadyOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return ErrOrderNotFound
	}
	if status == StatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if status == StatusDelivered {
		return ErrOrderAlreadyDelivered
	}
	if status == StatusReady {
		return ErrOrderAlreadyReady
	}

	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = StatusReady
	db.Save(order)

	return nil
}

// DeliverOrder marks the given order as deliered. Returns an error if the order does not exist or is already delivered or cancelled.
func DeliverOrder(id uint) error {
	status, err := GetOrderStatus(id)
	if err != nil {
		return ErrOrderNotFound
	}
	if status == StatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if status == StatusDelivered {
		return ErrOrderAlreadyDelivered
	}

	order := &models.Order{ID: id}
	db.Model(order).First(&order)
	order.Status = StatusDelivered
	db.Save(order)

	return nil
}
