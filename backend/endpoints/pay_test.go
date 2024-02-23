//go:build integration

package endpoints

import (
	"fmt"
	"net/http"
	"teamproject/database"
	"teamproject/database/models"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPay(t *testing.T) {
	app := fiber.New()
	app.Patch("/pay/:id", Pay)

	database.ClearMenu()
	database.ClearOrders()

	// Create test orders
	menuItem1 := models.MenuItem{Name: "Tequila"}
	menuItem2 := models.MenuItem{Name: "Vodka"}
	menuItem3 := models.MenuItem{Name: "Rum"}
	database.AddItem(&menuItem1)
	database.AddItem(&menuItem2)
	database.AddItem(&menuItem3)

	var testItemID uint = 1
	var testItemID2 uint = 2

	testItemList1 := []models.OrderItem{{Item: menuItem1, Notes: "Item1"}, {Item: menuItem2, Notes: "Notes2"}}
	testItemList2 := []models.OrderItem{{Item: menuItem3, Notes: "Notes3"}}

	testOrder := models.Order{ID: testItemID, Time: time.Now(), TableNumber: 16, Bill: 16.99, Status: "Ready", Items: testItemList1}
	testOrder2 := models.Order{ID: testItemID2, Time: time.Now(), TableNumber: 17, Bill: 17.99, Status: "Preparing", Items: testItemList2}

	database.AddOrder(&testOrder)
	database.AddOrder(&testOrder2)

	testCases := []struct {
		name string
		id   string
		code int
	}{
		{
			name: "WithValidID",
			id:   "1",
			code: fiber.StatusOK,
		},
		{
			name: "WithDuplicateID",
			id:   "1",
			code: fiber.StatusConflict,
		},
		{
			name: "WithSecondValidID",
			id:   "2",
			code: fiber.StatusOK,
		},
		{
			name: "WithInvalidID",
			id:   "3",
			code: fiber.StatusNotFound,
		},
		{
			name: "WithNoID",
			id:   "",
			code: fiber.StatusNotFound,
		},
		{
			name: "WithStringID",
			id:   "ABC",
			code: fiber.StatusUnprocessableEntity,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/pay/%s", test.id), nil)

			// Send test HTTP request
			res, err := app.Test(req)
			assert.NoError(t, err)
			defer res.Body.Close()

			// Check response code
			assert.Equal(t, test.code, res.StatusCode, "Test that the correct status code was returned")
		})
	}
}
