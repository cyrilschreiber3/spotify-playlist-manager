package database

import (
	"log"
	"os"
)

func Init() {
	log.Println("Initializing database...")
	if os.Getenv("DB_TYPE") == "mongodb" {
		mongoInit()
	}
}

func Close() {
	log.Println("Closing database connection...")
	if os.Getenv("DB_TYPE") == "mongodb" {
		mongoClose()
	}
}
