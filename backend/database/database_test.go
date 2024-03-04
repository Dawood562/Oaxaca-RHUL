//go:build integration

package database

import (
	"strconv"
	"teamproject/database/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseQueries(t *testing.T) {
	ResetTestMenu()

	testCases := []struct {
		name            string
		filter          *MenuFilter
		expectedLen     int
		expectedElement models.MenuItem
	}{
		{
			name:            "EmptyFilter",
			filter:          &MenuFilter{},
			expectedLen:     4,
			expectedElement: models.MenuItem{ID: 1, Name: "TESTFOOD", Description: "Description for TESTFOOD", Price: 5.00, Calories: 400, Allergens: []models.Allergen{}},
		},
		{
			name:            "WithSearchTermFilter",
			filter:          &MenuFilter{SearchTerm: "TESTFOOD2"},
			expectedLen:     1,
			expectedElement: models.MenuItem{ID: 2, Name: "TESTFOOD2", Description: "Description for TESTFOOD2", Price: 6.00, Calories: 500, Allergens: []models.Allergen{}},
		},
		{
			name:            "WithMultipleFilters",
			filter:          &MenuFilter{MaxPrice: 6.00, MaxCalories: 600},
			expectedLen:     2,
			expectedElement: models.MenuItem{ID: 1, Name: "TESTFOOD", Description: "Description for TESTFOOD", Price: 5.00, Calories: 400, Allergens: []models.Allergen{}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := QueryMenu(test.filter)
			assert.Equal(t, test.expectedLen, len(result), "Check query returns correct number of items")
			assert.Contains(t, result, test.expectedElement, "Check that query result contains expected item")
		})
	}
}

func TestDatabaseInserts(t *testing.T) {
	ClearMenu()
	item := &models.MenuItem{
		Name:     "TestInsert",
		Price:    5.00,
		Calories: 500,
	}
	err := AddItem(item)
	assert.NoError(t, err, "Test that inserting a record does not create an error")
	// Add again to check duplicate prevention
	err = AddItem(item)
	assert.Error(t, err, "Test that inserting a duplicate record creates an error")
	// Check the size of the database
	menu := QueryMenu(&MenuFilter{})
	assert.Equal(t, 1, len(menu), "Check that the record was added to the menu")

	// Add a different item
	item = &models.MenuItem{
		Name:     "TestInsert2",
		Price:    6.00,
		Calories: 600,
	}
	err = AddItem(item)
	assert.NoError(t, err, "Test that adding an item does not create an error")
	menu = QueryMenu(&MenuFilter{})
	assert.Equal(t, 2, len(menu), "Check that the second record was added to the menu")
}

func TestDatabaseDelete(t *testing.T) {
	ResetTestMenu()

	// Delete TESTFOOD4
	err := RemoveItem(4)
	assert.NoError(t, err, "Test that removing a record does not create an error")
	// Check the record was really removed
	menu := QueryMenu(&MenuFilter{})
	assert.Equal(t, 3, len(menu), "Check record was removed from menu")
	// Check removing a duplicate item
	err = RemoveItem(4)
	assert.Error(t, err, "Check that removing a non-existent item creates an error")
	menu = QueryMenu(&MenuFilter{})
	assert.Equal(t, 3, len(menu), "Check that no items were removed from the database")
}

func TestDatabaseEdit(t *testing.T) {
	ResetTestMenu()

	// Check that a valid record can be edited
	newItem := models.MenuItem{ID: 1, Name: "TESTFOOD5", Price: 6.00, Calories: 500, Allergens: []models.Allergen{}}
	err := EditItem(&newItem)
	assert.NoError(t, err, "Test that editing a valid record does not create an error")
	// Check that the fields were modified
	menu := QueryMenu(&MenuFilter{})
	assert.Contains(t, menu, newItem, "Test that the item was successfully edited")

	// Check that an invalid record can't be edited
	newItem = models.MenuItem{ID: 5, Name: "TESTFOOD3"}
	err = EditItem(&newItem)
	assert.Error(t, err, "Test that editing an invalid record creates an error")

	newItem = models.MenuItem{ID: 4, Name: "TESTFOOD2"}
	err = EditItem(&newItem)
	assert.Error(t, err, "Test that editing an item with a duplicate name creates an error")
}

func TestDBAuth(t *testing.T) {
	user, name, pass := fetchDBAuth()
	if user == "-1" || name == "-1" || pass == "-1" {
		t.Fail()
	}
}

func TestPrepareArgsEmpty(t *testing.T) {
	args := prepareArgs(&MenuFilter{})
	assert.Equal(t, "%", args.SearchTerm, "Test search term default")
	assert.Equal(t, float32(9999), args.MaxPrice, "Test price default")
	assert.Equal(t, 9999, args.MaxCalories, "Test calorie default")
}

func TestPrepareArgsNotEmpty(t *testing.T) {
	args := prepareArgs(&MenuFilter{
		SearchTerm:  "test",
		MaxPrice:    5.00,
		MaxCalories: 500,
	})

	assert.Equal(t, "%test%", args.SearchTerm, "Test search term preparation")
	assert.Equal(t, float32(5.00), args.MaxPrice, "Test price preparation")
	assert.Equal(t, 500, args.MaxCalories, "Test calorie preparation")
}

func TestOrderRetrievalUnfiltered(t *testing.T) {
	// Check no data is returned to when no orders are in table
	ClearOrders()
	data, err := FetchOrders(false, -1)
	assert.Equal(t, 0, len(data), "Length of orders received should be 0 when no orders placed. Instead, received: "+strconv.Itoa(len(data)))

	// Check example data is inserted into table and returned correctly
	testData1 := models.Order{ID: 1, Time: time.Now(), TableNumber: 69, Bill: 42, Status: "Unknown"}
	err = AddOrder(&testData1)
	assert.NoError(t, err, "No error should be returned when adding empty order to Order table")
	data, err = FetchOrders(false, -1)
	assert.Equal(t, 1, len(data), "Size of data is actually: "+strconv.Itoa(len(data)))
	assert.Equal(t, 1, int(data[0].ID), "Returned id does not match expected id")
	assert.Equal(t, 69, int(data[0].TableNumber), "Returned table number doesnt match expected")
	assert.Equal(t, float64(42), data[0].Bill, "Returned bill does not match expected")
	assert.Equal(t, "Unknown", data[0].Status, "Status of recieved is: "+data[0].Status)

	RemoveOrder(testData1.ID)
}

func TestOrderQueryUnfiltered(t *testing.T) {
	ClearOrders()
	items := ResetTestMenu()

	testItemList := []models.OrderItem{{Item: items[0]}, {Item: items[1]}, {Item: items[2]}}
	testOrder := models.Order{ID: 1, Time: time.Now(), TableNumber: 16, Bill: 16.99, Status: "Ready", Items: testItemList}

	err := AddOrder(&testOrder)
	if err != nil {
		t.Fail()
	}

	returnedData, err := FetchOrders(false, -1)
	assert.Equal(t, 1, len(returnedData), "Incorrect number of data returned")

	err = RemoveOrder(testOrder.ID)
	assert.NoError(t, err, "Shouldnt throw error on removal")
}

func TestOrderRetrievalRejectDuplicate(t *testing.T) {
	ClearOrders()
	// Add test data
	testData1 := models.Order{ID: 1, Time: time.Now(), TableNumber: 69, Bill: 42, Status: "Unknown"}
	AddOrder(&testData1)

	err := AddOrder(&testData1)
	assert.Error(t, err, "Error should be thrown when duplicate item is added to order table")

	// Remove test data
	RemoveOrder(testData1.ID)
}

func TestRemoveItem(t *testing.T) {
	// Add empty test item
	testData1 := models.Order{ID: 1, Time: time.Now(), TableNumber: 69, Bill: 42, Status: "Unknown"}
	AddOrder(&testData1)
	err := AddOrder(&testData1)
	assert.Error(t, err, "Error should be thrown when duplicate item is added to order table")

	data, err := FetchOrders(false, -1)
	assert.Equal(t, 1, len(data), "Test item was not successfully added correctly")

	// Attempt to remove test item normally
	err = RemoveOrder(testData1.ID)
	assert.NoError(t, err, "Item was not removed!")

	err = RemoveOrder(testData1.ID)
	assert.Error(t, err, "Attempting to remove an item that doesnt exist should throw an error")

	// Check if item was really removed
	data, err = FetchOrders(false, -1)
	assert.Equal(t, 0, len(data), "Numbers of orders remaining should be 0 after removing only test order")
}

func TestRemovingMultipleData(t *testing.T) {
	ResetTestOrders()

	data, err := FetchOrders(false, -1)

	assert.Equal(t, 2, len(data), "Order table should contain only 1 item")

	err = RemoveOrder(1)
	assert.NoError(t, err, "Removing item returned an error")

	data, err = FetchOrders(false, -1)
	assert.Equal(t, 1, len(data), "Order table should contain no items!")
	assert.Equal(t, uint(2), data[0].ID, "Incorrect item deleted!")

	RemoveOrder(2)
	data, err = FetchOrders(false, -1)
	assert.Equal(t, 0, len(data), "Order table should contain no items!")
}

func TestFetchingOrdersCorrectlyBringsOrderItems(t *testing.T) {
	ResetTestOrders()

	testData, err := FetchOrders(false, -1)

	assert.NoError(t, err, "Shouldnt throw an error here")
	assert.Equal(t, 2, len(testData), "Incorrect amount of data fetched")
	assert.Equal(t, "Item1", testData[0].Items[0].Notes, "Incorrect item retrieved from order fetch")
}

func TestPayOrder(t *testing.T) {
	ResetTestOrders()

	paid, err := OrderPaid(1)
	assert.NoError(t, err, "Test that retrieving payment status for a valid order does not create an error")
	assert.False(t, paid, "Test that fresh order is unpaid")

	paid, err = OrderPaid(4)
	assert.Error(t, err, "Test that retrieving payment status for an invalid order creates an error")
	assert.False(t, paid, "Test that an order that doesn't exist is not paid")

	err = PayOrder(1)
	assert.NoError(t, err, "Test that paying for valid order creates no error")
	err = PayOrder(1)
	assert.Error(t, err, "Test that paying for a paid order creates an error")
	err = PayOrder(3)
	assert.Error(t, err, "Test that paying for an invalid order creates an error")

	paid, err = OrderPaid(1)
	assert.NoError(t, err)
	assert.True(t, paid, "Test that payment status was updated")

	paid, err = OrderPaid(2)
	assert.NoError(t, err)
	assert.False(t, paid, "Test that payment status was unchanged")

	err = PayOrder(2)
	assert.NoError(t, err, "Test that paying for a second valid order creates no errors")
	paid, err = OrderPaid(2)
	assert.NoError(t, err)
	assert.True(t, paid, "Test that payment status was updated")
}

func TestConfirmOrder(t *testing.T) {
	ResetTestOrders()

	// Test that order status begins as expected
	status, err := GetOrderStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, "Awaiting Confirmation", status, "Test that orders begin with the correct status")
	_, err = GetOrderStatus(3)
	assert.Error(t, err, "Test that getting the status of an order that does not exist creates an error")

	err = ConfirmOrder(1)
	assert.NoError(t, err, "Test that no error is returned when confirming a valid order")
	err = ConfirmOrder(3)
	assert.Error(t, err, "Test that an error is returned when confirming an invalid order")
	status, err = GetOrderStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, "Preparing", status, "Test that order status was updated")
	status, err = GetOrderStatus(2)
	assert.NoError(t, err)
	assert.Equal(t, "Awaiting Confirmation", status)
	err = ConfirmOrder(1)
	assert.Error(t, err, "Test that confirming a confirmed order creates an error")
	err = ConfirmOrder(2)
	assert.NoError(t, err)
	status, err = GetOrderStatus(2)
	assert.NoError(t, err)
	assert.Equal(t, "Preparing", status)
}

func TestCancelOrder(t *testing.T) {
	ResetTestOrders()

	status, err := GetOrderStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, StatusAwaitingConfirmation, status, "Check that order begins with correct status")

	err = CancelOrder(1)
	assert.NoError(t, err, "Check that cancelling a valid order creates no error")
	err = CancelOrder(1)
	assert.Error(t, err, "Check that cancelling an order twice creates an error")
	err = CancelOrder(3)
	assert.Error(t, err, "Check that cancelling an invalid order creates an error")
	status, err = GetOrderStatus(1)
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, status, "Test that order status was correctly updated")
	status, err = GetOrderStatus(2)
	assert.NoError(t, err)
	assert.Equal(t, StatusAwaitingConfirmation, status)
	err = CancelOrder(2)
	assert.NoError(t, err, "Test that cancelling a second vaid order creates no errors")
	status, err = GetOrderStatus(2)
	assert.NoError(t, err)
	assert.Equal(t, StatusCancelled, status, "Test that the second order status was correctly updated")
}

func TestAddAndRetrieveAllergens(t *testing.T) {
	ClearMenu()
	ClearOrders()

	err := AddItem(&models.MenuItem{ID: 1, Name: "Test 1", Allergens: []models.Allergen{{Name: "Gluten"}, {Name: "Nuts"}}})
	assert.NoError(t, err, "Test that adding allergens to menu does not create errors")
	// Test allergen length
	var count int64
	db.Model(&models.Allergen{}).Where("1=1").Count(&count)
	assert.Equal(t, int64(2), count, "Test that the correct number of allergens are present in the database")

	err = AddItem(&models.MenuItem{ID: 2, Name: "Test 2"})
	assert.NoError(t, err, "Test that adding an order without allergens does not create errors")
	db.Model(&models.Allergen{}).Where("1=1").Count(&count)
	assert.Equal(t, int64(2), count, "Test that the correct number of allergens are present in the database")

	orders := QueryMenu(&MenuFilter{})
	assert.Equal(t, 2, len(orders), "Test that the correct number of menu items are returned")
	assert.Equal(t, 2, len(orders[0].Allergens), "Test that the order items contain the correct number of allergens")
	assert.Equal(t, 0, len(orders[1].Allergens), "Test that the order items contain the correct number of allergens")

	err = RemoveItem(1)
	assert.NoError(t, err, "Test that removing a menu item with allergens does not create an error")
	db.Model(&models.Allergen{}).Where("1=1").Count(&count)
	assert.Equal(t, int64(0), count, "Test that the correct number of allergens are present in the database")
}
