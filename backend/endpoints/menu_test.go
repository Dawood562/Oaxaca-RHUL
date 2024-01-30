//go:build integration

package endpoints

import (
	"testing"
)

func TestFetchMenu1(t *testing.T) {
	// Fetch all items on the menu
}

func TestFetchMenu2(t *testing.T) {
	// Fetch all items that have x calories
}

func TestFetchMenu3(t *testing.T) {
	// Fetch all items that do NOT have x allergen
}



func TestChangeMenu1(t *testing.T) {
	// Change item with given name to have a slightly different name, should return 200 Okay
}
func TestChangeMenu2(t *testing.T) {
	// Change name item with given ID, should return 200 Okay if exists and 404 if not
}
func TestChangeMenu3(t *testing.T) {
	// Change item with given ID to have different ID, should return 403
}