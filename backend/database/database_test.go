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
			expectedElement: models.MenuItem{ID: 1, Name: "TESTFOOD", Description: "Description for TESTFOOD", Price: 5.00, Calories: 400},
		},
		{
			name:            "WithSearchTermFilter",
			filter:          &MenuFilter{SearchTerm: "TESTFOOD2"},
			expectedLen:     1,
			expectedElement: models.MenuItem{ID: 2, Name: "TESTFOOD2", Description: "Description for TESTFOOD2", Price: 6.00, Calories: 500},
		},
		{
			name:            "WithMultipleFilters",
			filter:          &MenuFilter{MaxPrice: 6.00, MaxCalories: 600},
			expectedLen:     2,
			expectedElement: models.MenuItem{ID: 1, Name: "TESTFOOD", Description: "Description for TESTFOOD", Price: 5.00, Calories: 400},
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
	newItem := models.MenuItem{ID: 1, Name: "TESTFOOD5", Price: 6.00, Calories: 500}
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
	data := fetchOrders()
	assert.Equal(t, 0, len(data), "Length of orders received should be 0 when no orders placed. Instead, received: "+strconv.Itoa(len(data)))

	// Check example data is inserted into table and returned correctly
	testData1 := models.Order{ID: 1, Time: time.Now(), TableNumber: 69, Bill: 42, Status: "Unknown"}
	err := AddOrder(&testData1)
	assert.NoError(t, err, "No error should be returned when adding empty order to Order table")
	data = fetchOrders()
	assert.Equal(t, 1, len(data), "Size of data is actually: "+strconv.Itoa(len(data)))
	assert.Equal(t, 1, int(data[0].ID), "Returned id does not match expected id")
	assert.Equal(t, 69, int(data[0].TableNumber), "Returned table number doesnt match expected")
	assert.Equal(t, float64(42), data[0].Bill, "Returned bill does not match expected")
	assert.Equal(t, "Unknown", data[0].Status, "Status of recieved is: "+data[0].Status)

	RemoveOrder(testData1.ID)
}

func TestOrderQueryUnfiltered(t *testing.T) {
	menuItem1 := models.MenuItem{Name: "Tequila"}
	menuItem2 := models.MenuItem{Name: "Vodka"}
	menuItem3 := models.MenuItem{Name: "Rum"}

	var TestId1 uint = 1
	testItemList := []models.OrderItem{{OrderID: TestId1, Item: menuItem1}, {Item: menuItem2}, {Item: menuItem3}}
	testOrder := models.Order{ID: TestId1, Time: time.Now(), TableNumber: 16, Bill: 16.99, Status: "Ready", Items: testItemList}

	err := AddOrder(&testOrder)
	if err != nil {
		t.Fail()
	}

	returnedData := fetchOrders()
	assert.Equal(t, 1, len(returnedData), "Incorrect number of data returned")

	err = RemoveOrder(testOrder.ID)
	assert.NoError(t, err, "Shouldnt throw error on removal")
}

func TestOrderRetrievalRejectDuplicate(t *testing.T) {
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

	data := fetchOrders()
	assert.Equal(t, 1, len(data), "Test item was not successfully added correctly")

	// Attempt to remove test item normally
	err = RemoveOrder(testData1.ID)
	assert.NoError(t, err, "Item was not removed!")

	err = RemoveOrder(testData1.ID)
	assert.Error(t, err, "Attempting to remove an item that doesnt exist should throw an error")

	// Check if item was really removed
	data = fetchOrders()
	assert.Equal(t, 0, len(data), "Numbers of orders remaining should be 0 after removing only test order")
}

func TestRemovingMultipleData(t *testing.T) {
	menuItem1 := models.MenuItem{Name: "Tequila"}
	menuItem2 := models.MenuItem{Name: "Vodka"}
	menuItem3 := models.MenuItem{Name: "Rum"}

	var testItemID uint = 1
	var testItemID2 uint = 2

	testItemList1 := []models.OrderItem{{OrderID: testItemID, Item: menuItem1, Notes: "Item1"}, {OrderID: testItemID, Item: menuItem2, Notes: "Notes2"}}
	testItemList2 := []models.OrderItem{{OrderID: testItemID2, Item: menuItem3, Notes: "Notes3"}}

	testOrder := models.Order{ID: testItemID, Time: time.Now(), TableNumber: 16, Bill: 16.99, Status: "Ready", Items: testItemList1}
	testOrder2 := models.Order{ID: testItemID2, Time: time.Now(), TableNumber: 17, Bill: 17.99, Status: "Preparing", Items: testItemList2}

	AddOrder(&testOrder)
	AddOrder(&testOrder2)
	data := fetchOrders()

	assert.Equal(t, 2, len(data), "Order table should contain only 1 item")

	err := RemoveOrder(testItemID)
	if err != nil {
		assert.NoError(t, err, "Removing item returned an error")
	}

	data = fetchOrders()
	assert.Equal(t, 1, len(data), "Order table should contain no items!")

	RemoveOrder(testItemID2)
	data = fetchOrders()
	assert.Equal(t, 0, len(data), "Order table should contain no items!")
}
