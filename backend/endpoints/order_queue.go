package endpoints

import (
	"errors"
	"teamproject/database/models"
)

var orderQueue []models.Order

// Returns true if any orders have been placed on queue
func IsAnyPendingOrders() bool {
	return len(orderQueue) > 0
}

// Adds provided order to queue to be given to a waiter
func AddOrderToQueue(order models.Order) error {
	existing := false
	for _, o := range orderQueue {
		if order.ID == o.ID {
			existing = true
		}
	}
	if !existing {
		orderQueue = append(orderQueue, order)
		return nil
	} else {
		return errors.New("Order with same id already in queue")
	}
}

func clearOrderQueue() {
	orderQueue = []models.Order{}
}
