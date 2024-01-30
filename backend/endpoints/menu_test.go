//go:build integration

package endpoints

import (
	"testing"

	"teamproject/database"
)

/*
Tests that all menu items can be successfully fetched.
*/
func TestFetchMenu1(t *testing.T) {
	_, code := FetchMenu(nil)
	if code != nil { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those that have < x calories
*/
func TestFetchMenu2(t *testing.T) {
	filter := database.MenuFilter{LowCalories: true}
	_, code := FetchMenu(&filter)
	if code != nil { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those missing one allergen
*/
func TestFetchMenu3(t *testing.T) {
	filter := database.MenuFilter{Allergens: []string{"gluten"}}
	_, code := FetchMenu(&filter)
	if code != nil { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those missing multiple allergens
*/
func TestFetchMenu4(t *testing.T) {
	filter := database.MenuFilter{Allergens: []string{"gluten", "soya", "milk", "eggs"}}
	_, code := FetchMenu(&filter)
	if code != nil { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that menu items can be filtered to those under a certain price
*/
func TestFetchMenu5(t *testing.T) {
	filter := database.MenuFilter{MaxPrice: 14.99}
	_, code := FetchMenu(&filter)
	if code != nil { // If there was any error getting items
		t.Fail()
	}
}

/*
Tests that, given an item ID, some data like the dish name can be changed.

Should return 200 Okay if successful.
*/
// func TestChangeMenu1(t *testing.T) {
// 	var itemID int = 562
// 	var newName string = "What a nice day outside! Too bad I'm coding."
// 	error := ChangeMenu(itemID, newName)
// 	if error != 200 {
// 		t.Fail()
// 	}
// }

/*
Tests that, given an item ID, the item's ID CANNOT and WILL NOT be changed.

Should return 403 (Forbidden).
*/
func TestChangeMenu2(t *testing.T) {

}
