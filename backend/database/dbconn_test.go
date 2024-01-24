package database

import (
	"testing"
)

func TestUpdateDB(t *testing.T) {
	if !UpdateDB("INSERT INTO Customers VALUES (1, 'John', 'Dont kick his dog')") {
		t.Fail()
	}
}

func TestQueryDB(t *testing.T) {
	retrieved := QueryDB("SELECT * FROM customers")
	if retrieved.CustomerID == -1 {
		t.Fail()
	}
}
