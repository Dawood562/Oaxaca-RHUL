//go:build integration

package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseQueries(t *testing.T) {
	// Setup test data
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (1, 'TESTFOOD', 5.00, 400)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (2, 'TESTFOOD2', 6.00, 500)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (3, 'TESTFOOD3', 7.00, 600)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (4, 'TESTFOOD4', 8.01, 720)")

	testCases := []struct {
		name            string
		filter          *MenuFilter
		expectedLen     int
		expectedElement MenuItem
	}{
		{
			name:            "EmptyFilter",
			filter:          &MenuFilter{},
			expectedLen:     4,
			expectedElement: MenuItem{ID: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400},
		},
		{
			name:            "WithSearchTermFilter",
			filter:          &MenuFilter{SearchTerm: "TESTFOOD2"},
			expectedLen:     1,
			expectedElement: MenuItem{ID: 2, ItemName: "TESTFOOD2", Price: 6.00, Calories: 500},
		},
		{
			name:            "WithMultipleFilters",
			filter:          &MenuFilter{MaxPrice: 6.00, MaxCalories: 600},
			expectedLen:     2,
			expectedElement: MenuItem{ID: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := QueryMenu(test.filter)
			assert.Equal(t, test.expectedLen, len(result), "Check query returns correct number of items")
			assert.Contains(t, result, test.expectedElement, "Check that query result contains expected item")
		})
	}

	UpdateDB("DELETE FROM menuitem")
}

func TestDatabaseInserts(t *testing.T) {
	item := &MenuItem{
		ItemName: "TestInsert",
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
	item = &MenuItem{
		ItemName: "TestInsert2",
		Price:    6.00,
		Calories: 600,
	}
	err = AddItem(item)
	assert.NoError(t, err, "Test that adding an item does not create an error")
	menu = QueryMenu(&MenuFilter{})
	assert.Equal(t, 2, len(menu), "Check that the second record was added to the menu")

	UpdateDB("DELETE FROM menuitem")
}

func TestDatabaseDelete(t *testing.T) {
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (1, 'TESTFOOD', 5.00, 400)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (2, 'TESTFOOD2', 6.00, 500)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (3, 'TESTFOOD3', 7.00, 600)")
	UpdateDB("INSERT INTO menuitem (menuitemid, menuitemname, price, calories) VALUES (4, 'TESTFOOD4', 8.01, 720)")

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

	UpdateDB("DELETE FROM menuitem")
}

func TestDatabaseEdit(t *testing.T) {
	UpdateDB("INSERT INTO menuitem (menuItemId, menuItemName, price, calories) VALUES (1, 'TESTFOOD', 5.00, 400)")

	// Check that a valid record can be edited
	newItem := MenuItem{ID: 1, ItemName: "TESTFOOD2", Price: 6.00, Calories: 500}
	err := EditItem(&newItem)
	assert.NoError(t, err, "Test that editing a valid record does not create an error")
	// Check that the fields were modified
	menu := QueryMenu(&MenuFilter{})
	assert.Contains(t, menu, newItem, "Test that the item was successfully edited")

	// Check that an invalid record can't be edited
	newItem = MenuItem{ID: 2, ItemName: "TESTFOOD3"}
	err = EditItem(&newItem)
	assert.Error(t, err, "Test that editing an invalid record creates an error")

	UpdateDB("DELETE FROM menuitem")
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
