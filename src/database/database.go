// Package database provides functions for interacting with the PostgreSQL database. This package uses structs from the models package to represent data. The
// endpoints package is expected to provide data in this format, and this package will process it into an SQL query appropriate for the interaction. Data is
// returned in structs from the models package. Errors are returned by functions in this package if there is a database error, or if the provided data causes
// a database conflict.
package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"teamproject/database/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Globally available database instance
var db *gorm.DB

var (
	ErrOrderNotFound         error = errors.New("order not found")
	ErrOrderAlreadyPaid      error = errors.New("order already paid for")
	ErrOrderAlreadyConfirmed error = errors.New("order already confirmed")
	ErrOrderAlreadyCancelled error = errors.New("order already cancelled")
	ErrOrderAlreadyReady     error = errors.New("order already ready")
	ErrOrderAlreadyDelivered error = errors.New("order already delivered")
)

const (
	StatusAwaitingConfirmation string = "Awaiting Confirmation"
	StatusPreparing            string = "Preparing"
	StatusReady                string = "Ready"
	StatusDelivered            string = "Delivered"
	StatusCancelled            string = "Cancelled"
)

func init() {
	// Fetch the database auth credentials
	dbUsername, dbName, dbPassword := fetchDBAuth()
	url := "postgres://" + dbUsername + ":" + dbPassword + "@db:5432/" + dbName
	// Connect to database and store reference to it
	conn, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	db = conn
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	// AutoMigrate all tables
	err = db.AutoMigrate(&models.MenuItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.OrderItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&models.Allergen{})
	if err != nil {
		log.Fatal(err)
	}
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
