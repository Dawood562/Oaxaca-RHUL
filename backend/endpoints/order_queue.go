package endpoints

import "teamproject/database/models"

var orderQueue []models.Order

// Returns true if any orders have been placed on queue
func IsAnyPendingOrders() bool {
	return len(orderQueue) > 0
}

// Adds provided order to queue to be given to a waiter
func AddOrderToQueue(order models.Order) {
	orderQueue = append(orderQueue, order)
}

func clearOrderQueue() {
	orderQueue = []models.Order{}
}
