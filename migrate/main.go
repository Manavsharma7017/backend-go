package main

import (
	"backend/database"
)

func main() {
	database.ConnectDB()
	database.AutoMigrate()
}
