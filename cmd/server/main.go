package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/asino-nelson/safiri-logistics/internal/database"
	"github.com/asino-nelson/safiri-logistics/internal/user"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Connect DB
	database.ConnectDB()

	repo := user.NewRepository(database.DB)

	u := &user.User{
		ID:           uuid.New().String(),
		Name:         "Nelson",
		Email:        "nelson@test.com",
		PasswordHash: "hashedpassword",
		Role:         "customer",
	}

	err = repo.CreateUser(context.Background(), u)
	if err != nil {
		log.Fatal("Error creating user:", err)
	}

	log.Println("✅ User created")

	foundUser, err := repo.GetUserByEmail(context.Background(), "nelson@test.com")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetched user:", foundUser.Email)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.Run(":8080")
}
