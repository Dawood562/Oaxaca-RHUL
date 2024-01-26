package database

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type MenuItem struct {
	MenuItemId string  `json:"menu_item_id"`
	ItemName   string  `json:"item_name"`
	Price      float64 `json:"price"`
	Calories   float64 `json:"calories"`
}

func init() {
	dbUsername, dbPassword := fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@localhost:5432/teamproject30"
	dbLocal, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	db = dbLocal
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
Takes in clause to retrieve specific items from menu table.
Example: "itemname = 'chicken korma'"
If entire table required then leave clause empty
Returns struct of example customer for now
Returns -1 in customerID if unable to access database
*/
func QueryMenu(clause ...string) []*MenuItem {
	dbTable := db.Table("menuitem")
	// If clause provided, use it
	if len(clause) > 0 {
		dbTable.Where(clause[0])
	}
	rows, err := dbTable.Rows()
	if err != nil {
		log.Fatal(err)
	}

	var data []*MenuItem
	for rows.Next() {
		var _menuitemid string
		var _itemname string
		var _price float64
		var _calories float64
		rows.Scan(&_menuitemid, &_itemname, &_price, &_calories)
		data = append(data, &MenuItem{MenuItemId: _menuitemid, ItemName: _itemname, Price: _price, Calories: _calories})
	}

	return data
}

/*
Fetches database login details from db_details.txt file
db_details.txt should be in database folder with following content structure:
<username>
<password>
*/
func fetchDBAuth() (string, string) {

	file, err := os.Open("db_details.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)

	reader.Scan()
	dbUsername := reader.Text()
	reader.Scan()
	dbPassword := reader.Text()
	return dbUsername, dbPassword
}
