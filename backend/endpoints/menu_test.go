//go:build integration

package endpoints

import (
	"teamproject/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMenu(t *testing.T) {
	// Setup test data
	database.UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD', 5.00, 400)")
	database.UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD2', 6.00, 500)")
	database.UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD3', 7.00, 600)")
	database.UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD4', 8.01, 720)")

	testCases := []struct {
		name   string
		filter *database.MenuFilter
	}{
		{
			name:   "EmptyFilter",
			filter: &database.MenuFilter{},
		},
		{
			name:   "WithSearchTermFilter",
			filter: &database.MenuFilter{SearchTerm: "2"},
		},
		{
			name:   "WithPriceFilter",
			filter: &database.MenuFilter{MaxPrice: 5.00},
		},
		{
			name:   "WithCalorieFilter",
			filter: &database.MenuFilter{MaxCalories: 500},
		},
		{
			name: "WithMultipleFilters",
			filter: &database.MenuFilter{
				SearchTerm:  "TESTFOOD",
				MaxPrice:    6.00,
				MaxCalories: 600,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := FetchMenu(test.filter)
			assert.NoError(t, err)
		})
	}

	database.UpdateDB("DELETE FROM menuitem")
}
