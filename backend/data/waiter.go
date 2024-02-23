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

func ClearWaiterList() {
	activeWaiters = []Waiter{}
}
