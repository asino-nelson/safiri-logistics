package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/asino-nelson/safiri-logistics/internal/database"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	database.ConnectDB()

	row := database.DB.QueryRow(context.Background(), "SELECT NOW()")

	var timeNow string
	err = row.Scan(&timeNow)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Time:", timeNow)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.Run(":8080")
}
