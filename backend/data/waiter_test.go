//go:build integration

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddWaiter(t *testing.T) {
	ClearWaiterList()
	tableNumbers := []uint{1, 2, 3}
	AddWaiter(Waiter{ID: 69, Username: "John", TableNumber: tableNumbers})
	assert.Equal(t, 1, len(activeWaiters), "Adding a waiter should increase number of waiters in list by 1")
}

func TestRemoveWaiter(t *testing.T) {
	addGenericTestData()
	assert.Equal(t, 2, len(activeWaiters), "Adding two waiters should return active waiter list of 2")
	RemoveWaiter(Waiter{ID: 1})
	assert.Equal(t, 1, len(activeWaiters), "Removing a waiter should return a active waiter list of 1")
	assert.Equal(t, 2, int(activeWaiters[0].ID), "Remaining waiter does not match expected waiter. Incorrect waiter removed!")
}

func TestGetAllWaiters(t *testing.T) {
	addGenericTestData()
	data := *GetWaiter()
	assert.Equal(t, 2, len(data), "Incorrect number of waiters returned from waiter list")
	assert.Equal(t, "John", data[0].Username, "Incorrect waiter data returned")
	assert.Equal(t, "James", data[1].Username, "Incorrect waiter data returned")
}

func TestGetWaitersWithFilter(t *testing.T) {
	addGenericTestData()
	data := *GetWaiter(Waiter{ID: 1})
	assert.Equal(t, 1, len(data), "Incorrect number of waiters returned from waiter list")
}

func TestGetWaitersWithInvalidFilter(t *testing.T) {
	addGenericTestData()
	data := GetWaiter(Waiter{ID: 3})
	assert.Nil(t, data, "Fetching data with invalid id should return nil")
}

func addGenericTestData() {
	ClearWaiterList()
	tableNumbers := []uint{1, 2, 3}
	waiter1 := Waiter{ID: 1, Username: "John", TableNumber: tableNumbers}
	waiter2 := Waiter{ID: 2, Username: "James", TableNumber: tableNumbers}
	AddWaiter(waiter1)
	AddWaiter(waiter2)
}
