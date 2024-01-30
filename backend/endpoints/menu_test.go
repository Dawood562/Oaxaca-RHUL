//go:build integration

package endpoints

import (
	"testing"
	"structs"
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
	filter := MenuFilter{low_calories: true}
	menu, code := FetchMenu(filter)
	if code != 200 { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those missing one allergen
*/
func TestFetchMenu3(t *testing.T) {
	filter := MenuFilter{allergens: ["gluten"]}
	menu, code := FetchMenu(filter)
	if code != 200 { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those missing multiple allergens
*/
func TestFetchMenu4(t *testing.T) {
	filter := MenuFilter{allergens: ["gluten", "soya", "milk", "eggs"]}
	menu, code := FetchMenu(filter)
	if code != 200 { // If there was any error getting items
		t.Fail()
	}
}


/*
Tests that menu items can be filtered to those under a certain price
*/
func TestFetchMenu5(t *testing.T) {
	filter := MenuFilter{max_price: 14.99}
	menu, code := FetchMenu(filter)
	if code != 200 { // If there was any error getting items
		t.Fail()
	}
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
	
}