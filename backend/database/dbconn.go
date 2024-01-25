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

// TESTING PURPOSES
type Customer struct {
	gorm.Model
	CustomerID   int
	Name         string
	Instructions string
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
	db.Exec(stmt)
	closeDB(db)
	return true
}

/*
Returns struct of example customer for now
Returns -1 in customerID if query unsuccessful
*/
func QueryDB(stmt string) Customer {
	db := openDB()
	if db == nil {
		return Customer{CustomerID: -1}
	}
	var data Customer
	db.Raw(stmt).Scan(&data)
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

	// USED IN TESTING - TO BE CHANGED WHEN ACTUAL DATABASE USED
	db.AutoMigrate(&Customer{})

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
		return "-1", "-1"
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
