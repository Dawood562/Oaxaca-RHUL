package endpoints

import "teamproject/database/models"

var orderQueue []models.Order

func IsAnyPendingOrders() bool {
	return len(orderQueue) > 0
}

func clearOrderQueue() {
	orderQueue = []models.Order{}
}
