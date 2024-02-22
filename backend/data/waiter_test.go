//go:build integration

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWaiter(t *testing.T) {
	tableNumbers := []uint{1, 2, 3}
	AddWaiter(Waiter{ID: 69, Username: "John", TableNumber: tableNumbers})
	assert.Equal(t, 1, len(ActiveWaiters), "Adding a waiter should increase number of waiters in list by 1")
}
