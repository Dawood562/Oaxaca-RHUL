package database

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbUsername string
var dbPassword string
var dbFetchedAuth bool

type MenuItem struct {
	MenuItemId string
	ItemName   string
	Price      float64
	Calories   float64
}

/*
Executes SQL query provided
returns true if query exeucted successfully
*/
func UpdateDB(stmt string) bool {
	db := openDB()
	if db == nil {
		return false
	}
	//db.Exec(stmt)
	closeDB(db)
	return true
}

/*
Takes in clause to retrieve specific items from menu table.
Example: "itemname = 'chicken korma'"
If entire table required then leave clause empty
Returns struct of example customer for now
Returns -1 in customerID if unable to access database
*/
func QueryMenu(clause ...string) []MenuItem {
	db := openDB()
	if db == nil {
		return []MenuItem{MenuItem{MenuItemId: "-1"}}
	}

	dbTable := db.Table("menuitem")
	// If clause provided, use it
	if len(clause) > 0 {
		dbTable.Where(clause[0])
	}
	rows, err := dbTable.Rows()
	if err != nil {
		log.Fatal(err)
	}

	var data []MenuItem
	for rows.Next() {
		var _menuitemid string
		var _itemname string
		var _price float64
		var _calories float64
		rows.Scan(&_menuitemid, &_itemname, &_price, &_calories)
		data = append(data, MenuItem{MenuItemId: _menuitemid, ItemName: _itemname, Price: _price, Calories: _calories})
	}

	closeDB(db)
	return data
}

func openDB() *gorm.DB {
	fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@localhost:5432/teamproject30"
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		fmt.Println("Successfully connected to database!")
	}

	return db
}

func closeDB(db *gorm.DB) {
	conn, _ := db.DB()
	conn.Close()
}

/*
Fetches database login details from db_details.txt file
db_details.txt should be in database folder with following content structure:
<username>
<password>
*/
func fetchDBAuth() (string, string) {
	if dbFetchedAuth {
		return dbUsername, dbPassword
	}

	file, err := os.Open("db_details.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)

	reader.Scan()
	dbUsername = reader.Text()
	reader.Scan()
	dbPassword = reader.Text()
	dbFetchedAuth = true
	return dbUsername, dbPassword
}
