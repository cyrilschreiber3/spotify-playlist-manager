package database

import "os"

func Init() {
	if os.Getenv("DB_TYPE") == "mongodb" {
		mongoInit()
	}
}

func Close() {
	if os.Getenv("DB_TYPE") == "mongodb" {
		mongoClose()
	}
}
