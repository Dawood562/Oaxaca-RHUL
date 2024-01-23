package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CustomerID   int
	Name         string
	Instructions string
}

func UpdateDB(stmt string) {
	db, err := gorm.Open(postgres.Open("postgres"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Customer{})
	db.Exec(stmt)
	conn, _ := db.DB()
	conn.Close()
}
