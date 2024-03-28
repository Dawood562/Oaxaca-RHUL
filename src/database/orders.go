package database

import (
	"errors"
	"fmt"
	"teamproject/database/models"
)

// AddOrder adds the given order to the database.
// Returns an error if an order exists with status != "Delivered" for the same table number.
func AddOrder(item *models.Order) error {
	// Check that there are no open orders with that table number already
	var count int64
	db.Model(&models.Order{}).Where("status NOT IN ? AND table_number = ?", []string{StatusDelivered, StatusCancelled}, item.TableNumber).Count(&count)
	if count > 0 {
		return fmt.Errorf("there is already an open order for table %d", item.TableNumber)
	}
	// Do not re-create menu items when adding them to an OrderItem
	feedback := db.Omit("Items.Item.*").Create(item)
	return feedback.Error
}

// RemoveOrder removes the order with the given ID from the database.
// Returns an error if no orders exist with that ID.
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

// FetchOrders returns a slice of orders.
// If confirmed is true, only returns orders which have been marked as confirmed.
// Returns an error if there is an error retrieving the data.
func FetchOrders(confirmed bool) ([]*models.Order, error) {
	dbCopy := *db
	dbLocal := &dbCopy
	var orderData []*models.Order

	if confirmed {
		dbLocal = dbLocal.Where("Status != ?", StatusAwaitingConfirmation)
	}
	result := dbLocal.Model(&models.Order{}).Preload("Items").Preload("Items.Item").Preload("Items.Item.Allergens").Find(&orderData)

	return orderData, result.Error
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
