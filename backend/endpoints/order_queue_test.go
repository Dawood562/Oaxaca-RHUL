package endpoints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPendingOrdersEmpty(t *testing.T) {
	clearOrderQueue()
	assert.False(t, IsAnyPendingOrders(), "Pending orders should return false when no items in pending order queue")
}
