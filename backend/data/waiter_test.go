//go:build integration

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWaiter(t *testing.T) {
	tableNumbers := []uint{1, 2, 3}
	AddWaiter(Waiter{ID: 69, Username: "John", TableNumber: tableNumbers})
	assert.Equal(t, 1, len(activeWaiters), "Adding a waiter should increase number of waiters in list by 1")
}

func TestRemoveWaiter(t *testing.T) {
	removePreviousWaiterData()
	tableNumbers := []uint{1, 2, 3}
	waiter1 := Waiter{ID: 1, Username: "John", TableNumber: tableNumbers}
	waiter2 := Waiter{ID: 2, Username: "James", TableNumber: tableNumbers}
	AddWaiter(waiter1)
	AddWaiter(waiter2)
	assert.Equal(t, 2, len(activeWaiters), "Adding two waiters should return active waiter list of 2")
	RemoveWaiter(Waiter{ID: 1})
	assert.Equal(t, 1, len(activeWaiters), "Removing a waiter should return a active waiter list of 1")
	assert.Equal(t, 2, int(activeWaiters[0].ID), "Remaining waiter does not match expected waiter. Incorrect waiter removed!")
}

func removePreviousWaiterData() {
	activeWaiters = []Waiter{} // Replace old waiter list with empty one
}
