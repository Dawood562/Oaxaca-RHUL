package endpoints

import (
	"teamproject/database/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPendingOrdersEmpty(t *testing.T) {
	clearOrderQueue()
	assert.False(t, IsAnyPendingOrders(), "Pending orders should return false when no items in pending order queue")
}

func TestOrderAddedCorrectly(t *testing.T) {
	clearOrderQueue()
	assert.False(t, IsAnyPendingOrders(), "Pending orders should return false when no items in pending order queue")
	testOrder1 := models.Order{ID: 1}
	AddOrderToQueue(testOrder1)
	assert.True(t, IsAnyPendingOrders(), "Pending orders should return true when an item is added to the order queue")
}
