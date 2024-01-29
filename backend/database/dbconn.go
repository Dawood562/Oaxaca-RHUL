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
	MenuItemId string  `json:"menu_item_id" gorm:"column:menuitemid"`
	ItemName   string  `json:"item_name" gorm:"column:itemname"`
	Price      float64 `json:"price" gorm:"column:price"`
	Calories   float64 `json:"calories" gorm:"column:calories"`
}

func init() {
	dbUsername, dbName, dbPassword := fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@db:5432/" + dbName
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
func QueryMenu(clause ...string) []MenuItem {

	var data []MenuItem
	db = db.Table("menuitem").Model(&MenuItem{})

	if len(clause) > 0 {
		db = db.Where(clause[0])
	}

	db.Find(&data)
	return data
}

/*
Fetches database login details from db_details.txt file
db_details.txt should be in database folder with following content structure:
<username>
<database name>
<password>
*/
func fetchDBAuth() (string, string, string) {
	file, err := os.Open("db_details.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewScanner(file)

	reader.Scan()
	dbUsername := reader.Text()
	reader.Scan()
	dbName := reader.Text()
	reader.Scan()
	dbPassword := reader.Text()
	return dbUsername, dbName, dbPassword
}
