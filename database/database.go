package database

import (
	"fmt"
	"log" // This helps us print messages in case something goes wrong.
	"messaging-platform/models"
	"os"

	"gorm.io/driver/postgres" // This helps us connect to a PostgreSQL database.
	"gorm.io/gorm"            // This is a tool that helps us work with our database.
)

// We have a special box called DB to keep our database connection.
var DB *gorm.DB

// ConnectDatabase is a function that helps us connect to our database.
func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	// This is like a secret code that tells us how to connect to our database.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	println(dsn)
	// We try to open a connection to our database using the secret code.
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// If something goes wrong, we print a message and stop the program.
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// We make sure our database knows about our data (like users, groups, messages, etc.).
	if err := database.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.GroupMembership{},
		&models.Message{},
		&models.MessageThread{},
		&models.Reaction{},
	); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
	// We put the database connection in our special box called DB.
	DB = database
}
