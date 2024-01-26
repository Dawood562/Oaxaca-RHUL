//go:build integration

package database

import (
	"testing"
)

func TestUpdateMenu(t *testing.T) {
	UpdateDB("DELETE FROM menuitem WHERE itemname = 'TESTFOOD';")
	UpdateDB("INSERT INTO menuitem VALUES ('999', 'TESTFOOD', 4.20, 450)")

	retrieved := QueryMenu("itemname='TESTFOOD'")[0]
	if retrieved.ItemName != "TESTFOOD" {
		t.Fail()
	}

	UpdateDB("DELETE FROM menuitem WHERE itemname = 'TESTFOOD';")
}

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
