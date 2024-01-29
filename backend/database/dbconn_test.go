//go:build integration

package database

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Always add TESTFOOD before testing
	UpdateDB("DELETE FROM menuitem WHERE itemname = 'TESTFOOD';")
	UpdateDB("INSERT INTO menuitem VALUES ('999', 'TESTFOOD', 4.20, 450)")

	code := m.Run()

	// Remove TESTFOOD after testing
	UpdateDB("DELETE FROM menuitem WHERE itemname = 'TESTFOOD';")

	os.Exit(code)
}

func TestUpdateMenu(t *testing.T) {
	retrieved := QueryMenu("itemname='TESTFOOD'")[0]
	if retrieved.ItemName != "TESTFOOD" {
		t.Fail()
	}
}

func TestUpdateQueryMultiple(t *testing.T) {
	// Add extra item to ensure its not retrieving incorrect item
	UpdateDB("INSERT INTO menuitem VALUES ('111', 'TESTFOOD2', 7.20, 720)")
	retrieved := QueryMenu("itemname='TESTFOOD'")
	if len(retrieved) != 1 {
		t.Fail() // if retrieved more than one item, should fail
	} else if retrieved[0].ItemName != "TESTFOOD" {
		t.Fail()
	}
	UpdateDB("DELETE FROM menuitem WHERE itemname='TESTFOOD2'")
}

func TestQueryMenuGeneric(t *testing.T) {
	retrieved := QueryMenu()
	if retrieved[0].MenuItemId == "-1" {
		t.Fail()
	}
}

func TestDBAuth(t *testing.T) {
	user, name, pass := fetchDBAuth()
	if user == "-1" || name == "-1" || pass == "-1" {
		t.Fail()
	}
}
