//go:build integration

package endpoints

import (
	"net/http"
	"teamproject/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRemoveItem(t *testing.T) {
	app := fiber.New()
	app.Delete("/remove_item", RemoveItem)

	database.AddItem(&database.MenuItem{
		ID:       1,
		Name:     "TESTFOOD",
		Price:    5.00,
		Calories: 400,
	})
	database.AddItem(&database.MenuItem{
		ID:       2,
		Name:     "TESTFOOD2",
		Price:    6.00,
		Calories: 500,
	})
	database.AddItem(&database.MenuItem{
		ID:       3,
		Name:     "TESTFOOD3",
		Price:    7.00,
		Calories: 600,
	})
	database.AddItem(&database.MenuItem{
		ID:       4,
		Name:     "TESTFOOD4",
		Price:    8.01,
		Calories: 720,
	})

	testCases := []struct {
		name              string
		id                string
		code              int
		expectedItemNames []string
	}{
		{
			name:              "WithCorrectFields",
			id:                "1",
			code:              200,
			expectedItemNames: []string{"TESTFOOD2", "TESTFOOD3", "TESTFOOD4"},
		},
		{
			name:              "WithEmptyId",
			id:                "",
			code:              422,
			expectedItemNames: []string{"TESTFOOD2", "TESTFOOD3", "TESTFOOD4"},
		},
		{
			name:              "WithInvalidId",
			id:                "1",
			code:              409,
			expectedItemNames: []string{"TESTFOOD2", "TESTFOOD3", "TESTFOOD4"},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("DELETE", "/remove_item", nil)
			q := req.URL.Query()
			q.Add("itemId", test.id)
			req.URL.RawQuery = q.Encode()

			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check the response
			assert.Equal(t, test.code, res.StatusCode, "Check that request returned expected status code")

			// Check that the database contains the correct items
			checkItemNames(t, test.expectedItemNames)
		})
	}

	database.ClearMenu()
}
