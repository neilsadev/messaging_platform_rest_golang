package main

import (
	"log"
	"messaging-platform/database"
	"messaging-platform/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()

	server := gin.Default()
	server.Use(cors.Default())

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
