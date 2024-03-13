package endpoints

import (
	"teamproject/database/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPendingOrdersEmpty(t *testing.T) {
	ClearOrderQueue()
	assert.False(t, IsAnyPendingOrders(), "Pending orders should return false when no items in pending order queue")
}

func TestOrderAddedCorrectly(t *testing.T) {
	ClearOrderQueue()
	assert.False(t, IsAnyPendingOrders(), "Pending orders should return false when no items in pending order queue")
	testOrder1 := models.Order{ID: 1}
	AddOrderToQueue(testOrder1)
	assert.True(t, IsAnyPendingOrders(), "Pending orders should return true when an item is added to the order queue")
}

func TestDuplicateOrderThrowsError(t *testing.T) {
	resetTestData(t)
}

func TestOrdersReturnCorrectly(t *testing.T) {
	resetTestData(t)
	queueData := GetQueuedOrders()
	assert.Equal(t, 2, len(queueData), "Returned incorrect number of orders")
	assert.Equal(t, "Ready", queueData[0].Status, "Incorrect data retrieved")
	assert.Equal(t, "Pending", queueData[1].Status, "Incorrect data retrieved")
}

func resetTestData(t *testing.T) {
	ClearOrderQueue()
	testOrder1 := models.Order{ID: 1, Status: "Ready"}
	testOrder2 := models.Order{ID: 1}
	testOrder3 := models.Order{ID: 2, Status: "Pending"}
	assert.NoError(t, AddOrderToQueue(testOrder1), "Adding an order shouldnt throw an error")
	assert.Error(t, AddOrderToQueue(testOrder2), "Adding a duplicate order should throw an error")
	assert.NoError(t, AddOrderToQueue(testOrder3), "Adding an order shouldnt throw an error")
}
