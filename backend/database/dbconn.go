package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dbUsername, dbName, dbPassword := fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@db:5432/" + dbName
	conn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	db = conn
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}
}

/*
Executes SQL query provided
*/
func UpdateDB(stmt string) {
	db.Exec(stmt)
}

/*
Takes in MenuFilter containing conditions to filter menu item.
If entire table required then provide empty MenuFilter
Returns array of MenuItems
*/
func QueryMenu(filter *MenuFilter) []MenuItem {
	preparedFilter := prepareArgs(filter)

	var data []MenuItem
	db.Table("menuitem").Model(&MenuItem{}).Where("itemname LIKE ?", preparedFilter.SearchTerm).Where("calories <= ?", preparedFilter.MaxCalories).Where("price <= ?", preparedFilter.MaxPrice).Find(&data)
	return data
}

// prepareArgs applies defaults to a MenuFilter struct in preparation for use in a query
func prepareArgs(filter *MenuFilter) *MenuFilter {
	ret := &MenuFilter{}

	// If the search term is shorter than 3 chars, disregard it
	if len(filter.SearchTerm) < 3 {
		ret.SearchTerm = "%"
	} else {
		ret.SearchTerm = "%" + filter.SearchTerm + "%"
	}

	// If no max calories are provided, set it to a high number
	if filter.MaxCalories <= 0 {
		ret.MaxCalories = 9999
	} else {
		ret.MaxCalories = filter.MaxCalories
	}

	// If no max price is provided, set it to a high number
	if filter.MaxPrice <= 0 {
		ret.MaxPrice = 9999
	} else {
		ret.MaxPrice = filter.MaxPrice
	}

	return ret
}

/*
Fetches database login details from .env
db_details.txt should be in database folder with following content structure:
<username>
<database name>
<password>
*/
func fetchDBAuth() (string, string, string) {
	err := godotenv.Load()
	if err != nil {
		// Load test env if production env not found
		godotenv.Load(".env.test")
	}
	username := os.Getenv("USERNAME")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	return username, dbname, password
}
