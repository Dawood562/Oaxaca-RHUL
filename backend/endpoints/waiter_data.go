package endpoints

import (
	"errors"
	"fmt"
	"strconv"
	"teamproject/database"
)

var activeWaiters []WaiterData

type WaiterData struct {
	ID          uint   `json:"id"`
	Username    string `json:"waiterUsername"`
	TableNumber []uint `json:"tableNumber"` // Array of table numbers that the waiter is attending
}

// Add waiter to current waiter list
// Returns true if waiter is successfully added to list
func AddWaiterData(waiter WaiterData) error {
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
func RemoveWaiterData(waiter WaiterData) error {
	for index, w := range activeWaiters {
		// If requested waiter is in active waiter list then remove it
		if waiter.ID == w.ID {
			ReAllocateTableNumbers(w)
			activeWaiters = append(activeWaiters[:index], activeWaiters[index+1:]...)

			return nil
		}
	}
	return errors.New("DID NOT FIND WAITER WITH PROVIDED ID TO REMOVE")
}

func ReAllocateTableNumbers(waiter WaiterData) error {
	if len(activeWaiters) <= 1 {
		fmt.Println("No active waiters")
		// If no waiters available add to queue

		// for each table number assigned to the waiter
		for _, wtn := range waiter.TableNumber {
			fmt.Println("Waiters table number:" + strconv.Itoa(int(wtn)))
			o, err := database.FetchOrders(false)
			if err != nil {
				return err
			}
			// Get orders and check where the assigned table number equals the fetched order table number
			for _, ord := range o {
				fmt.Println("Retrieved order table number:" + strconv.Itoa(int(ord.TableNumber)))
				if ord.TableNumber == wtn {
					fmt.Println("Added order to queue...")
					AddOrderToQueue(*ord)
				}
			}
		}
	} else {
		fmt.Println("Waiters active!")
	}
	return nil
}

// Gets waiter by provided id or if no waiter provided then gets all waiters
// Returns nil if waiter id provided and none
// Returns list of waiters
func GetWaiter(waiter ...WaiterData) *[]WaiterData {
	if len(waiter) > 0 {
		for _, w := range activeWaiters {
			if waiter[0].ID == w.ID {
				return &[]WaiterData{w}
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
	activeWaiters = []WaiterData{}
}
