package data

var activeWaiters []Waiter

type Waiter struct {
	ID          uint   `json:"id"`
	Username    string `json:"waiterUsername"`
	TableNumber []uint `json:"tableNumber"` // Array of table numbers that the waiter is attending
}

// Add waiter to current waiter list
func AddWaiter(waiter Waiter) {
	activeWaiters = append(activeWaiters, waiter)
}

// Remove waiter according to id provided in filter waiter
func RemoveWaiter(waiter Waiter) {
	for index, w := range activeWaiters {
		// If requested waiter is in active waiter list then remove it
		if waiter.ID == w.ID {
			activeWaiters = append(activeWaiters[:index], activeWaiters[index+1:]...)
		}
	}
}
