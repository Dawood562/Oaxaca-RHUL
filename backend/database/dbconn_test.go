package database

import (
	"testing"
)

// func TestUpdateDB(t *testing.T) {
// 	if !UpdateDB("INSERT INTO Customers VALUES (1, 'John', 'Dont kick his dog')") {
// 		t.Fail()
// 	}
// }

func TestQueryMenuGeneric(t *testing.T) {
	retrieved := QueryMenu()
	if retrieved[0].MenuItemId == "-1" {
		t.Fail()
	}
}

func TestQueryMenuWithClause(t *testing.T) {
	expectedExampleItemId := "1"
	expectedExampleItemName := "Food1"
	expectedExampleItemPrice := 3.99
	expectedExampleitemCalories := 250.0

	retrieved := QueryMenu("itemname = 'Food1'")[0]
	if retrieved.ItemName != expectedExampleItemName ||
		retrieved.MenuItemId != expectedExampleItemId ||
		retrieved.Price != expectedExampleItemPrice ||
		retrieved.Calories != expectedExampleitemCalories {
		t.Fail()
	}
}

func TestDBAuth(t *testing.T) {
	user, pass := fetchDBAuth()
	if user == "-1" || pass == "-1" {
		t.Fail()
	}
}
