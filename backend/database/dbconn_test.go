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

func TestQueryMenuGeneric(t *testing.T) {
	retrieved := QueryMenu()
	if retrieved[0].MenuItemId == "-1" {
		t.Fail()
	}
}

func TestDBAuth(t *testing.T) {
	user, pass := fetchDBAuth()
	if user == "-1" || pass == "-1" {
		t.Fail()
	}
}
