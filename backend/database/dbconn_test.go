//go:build integration

package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Setup test data
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD', 5.00, 400)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD2', 6.00, 500)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD3', 7.00, 600)")
	UpdateDB("INSERT INTO menuitem (itemname, price, calories) VALUES ('TESTFOOD4', 8.01, 720)")

	code := m.Run()

	// Remove test values after tests
	UpdateDB("DELETE FROM menuitem")

	os.Exit(code)
}

func TestGetMenu(t *testing.T) {
	menu := QueryMenu(&MenuFilter{})
	assert.Equal(t, 4, len(menu), "Number of menu items should be 4")
	assert.Contains(t, menu, MenuItem{MenuItemId: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400}, "Test that the correct menu items are returned")
}

func TestGetMenuWithArgs(t *testing.T) {
	items := QueryMenu(&MenuFilter{SearchTerm: "TESTFOOD2"})
	assert.Equal(t, 1, len(items), "Query should only return one item")
	assert.Equal(t, "TESTFOOD2", items[0].ItemName, "Test correct item is returned from query")
}

func TestGetMenuWithMultipleArgs(t *testing.T) {
	items := QueryMenu(&MenuFilter{MaxPrice: 6.00, MaxCalories: 600})
	assert.Equal(t, 2, len(items), "Query should return two items")
	assert.Contains(t, items, MenuItem{MenuItemId: 1, ItemName: "TESTFOOD", Price: 5.00, Calories: 400}, "Test that the correct menu items are returned")
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
