package main

import (
	"ki-d-assignment/config"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	
	db := config.SetupDatabaseConnection()
	config.CloseDatabaseConnection(db)

	router := gin.Default()
	router.Run()
}