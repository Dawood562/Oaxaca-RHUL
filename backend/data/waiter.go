package data

import (
	"errors"
	"strconv"
)

var activeWaiters []Waiter

type Waiter struct {
	ID          uint   `json:"id"`
	Username    string `json:"waiterUsername"`
	TableNumber []uint `json:"tableNumber"` // Array of table numbers that the waiter is attending
}

// Add waiter to current waiter list
// Returns true if waiter is successfully added to list
func AddWaiter(waiter Waiter) error {
	// Check that waiter with same id doesnt exist
	alreadyExists := false
	for _, existing := range activeWaiters {
		if existing.ID == waiter.ID {
			alreadyExists = true
		}
	}
	if !alreadyExists {
		activeWaiters = append(activeWaiters, waiter)
		return nil
	}
	return errors.New("Waiter cannot be added with duplicate id: " + strconv.Itoa(int(waiter.ID)))
}

// Remove waiter according to id provided in filter waiter
func RemoveWaiter(waiter Waiter) error {
	for index, w := range activeWaiters {
		// If requested waiter is in active waiter list then remove it
		if waiter.ID == w.ID {
			activeWaiters = append(activeWaiters[:index], activeWaiters[index+1:]...)
			return nil
		}
	}
	return errors.New("Did not find waiter with provided id to remove!")
}

// Gets waiter by provided id or if no waiter provided then gets all waiters
// Returns nil if waiter id provided and none
// Returns list of waiters
func GetWaiter(waiter ...Waiter) *[]Waiter {
	if len(waiter) > 0 {
		for _, w := range activeWaiters {
			if waiter[0].ID == w.ID {
				return &[]Waiter{w}
			}
		}
		return nil
	} else {
		return &activeWaiters
	}
}

func AddTableNumber(id uint, tableNumber uint) error {
	for index, waiter := range activeWaiters {
		if id == waiter.ID {
			activeWaiters[index].TableNumber = append(activeWaiters[index].TableNumber, tableNumber)
			return nil
		}
	}
	return errors.New("Waiter not found. Could not add waiter to table")
}

func ClearWaiterList() {
	activeWaiters = []Waiter{}
}
