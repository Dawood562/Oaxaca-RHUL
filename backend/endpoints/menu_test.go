//go:build integration

package endpoints

import (
	"teamproject/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMenu(t *testing.T) {
	filter := &database.MenuFilter{}
	_, err := FetchMenu(filter)
	assert.NoError(t, err)
}

func TestFetchMenuWithNameFilter(t *testing.T) {
	filter := &database.MenuFilter{
		SearchTerm: "2",
	}
	_, err := FetchMenu(filter)
	assert.NoError(t, err)
}

func TestFetchMenuWithPriceFilter(t *testing.T) {
	filter := &database.MenuFilter{
		MaxPrice: 5.00,
	}
	_, err := FetchMenu(filter)
	assert.NoError(t, err)
}

func TestFetchMenuWithCalorieFilter(t *testing.T) {
	filter := &database.MenuFilter{
		MaxCalories: 500,
	}
	_, err := FetchMenu(filter)
	assert.NoError(t, err)
}

func TestFetchMenuWithMultipleFilters(t *testing.T) {
	filter := &database.MenuFilter{
		SearchTerm:  "TESTFOOD",
		MaxPrice:    6.00,
		MaxCalories: 600,
	}
	_, err := FetchMenu(filter)
	assert.NoError(t, err)
}
