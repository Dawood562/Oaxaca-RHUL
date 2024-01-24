package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	return data
}

func openDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open("postgres://USERNAME_HERE:PASSWORD_HERE@localhost:PORT_HERE/DB_NAME_HERE"), &gorm.Config{})
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
