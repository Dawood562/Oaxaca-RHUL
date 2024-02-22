package data

var ActiveWaiters []Waiter

type Waiter struct {
	ID          uint   `json:"id"`
	Username    string `json:"waiterUsername"`
	TableNumber []uint `json:"tableNumber"` // Array of table numbers that the waiter is attending
}

func AddWaiter(waiter Waiter) {
	ActiveWaiters = append(ActiveWaiters, waiter)
}
