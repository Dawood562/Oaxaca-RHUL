//go:build integration

package endpoints

import (
	"testing"
)

/*
Tests that all menu items can be successfully fetched.
*/
func TestFetchMenu1(t *testing.T) {
	menu, code := FetchMenu()
	if code != 200 { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those that have < x calories
*/
func TestFetchMenu2(t *testing.T) {
	
}

/*
Tests that menu items can be filtered to those missing one or more allergens
*/
func TestFetchMenu3(t *testing.T) {

}

/*
Tests that menu items can be filtered to those under a certain price
*/
func TestFetchMenu4(t *testing.T) {
	// Fetch all items that under x price
}



/*
Tests that, given an item ID, some data like the dish name can be changed.

Should return 200 Okay if successful.
*/
func TestChangeMenu1(t *testing.T) {
	var itemID int = 562;
	var newName string = "What a nice day outside! Too bad I'm coding."
	error := ChangeMenu(itemID, newName)
	if error != 200 {
		t.Fail()
	}
}

/*
Tests that, given an item ID, the item's ID CANNOT and WILL NOT be changed.

Should return 403 (Forbidden).
*/

func TestChangeMenu2(t *testing.T) {
	// Change item with given ID to have different ID, should return 403
}