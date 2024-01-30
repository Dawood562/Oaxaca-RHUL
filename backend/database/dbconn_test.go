//go:build integration

package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseQueries(t *testing.T) {
	// Setup test data
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD', 5.00, 400)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD2', 6.00, 500)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD3', 7.00, 600)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD4', 8.01, 720)")

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
			expectedElement: MenuItem{MenuItemId: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400},
		},
		{
			name:            "WithSearchTermFilter",
			filter:          &MenuFilter{SearchTerm: "TESTFOOD2"},
			expectedLen:     1,
			expectedElement: MenuItem{MenuItemId: 2, ItemName: "TESTFOOD2", Price: 6.00, Calories: 500},
		},
		{
			name:            "WithMultipleFilters",
			filter:          &MenuFilter{MaxPrice: 6.00, MaxCalories: 600},
			expectedLen:     2,
			expectedElement: MenuItem{MenuItemId: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400},
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
