package main

import (
	"fmt"
	"log"
	"net/http"
	"user-crud-api/config"
	"user-crud-api/models"
	"user-crud-api/routes"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// Auto migrate User table
	config.DB.AutoMigrate(&models.User{})

	fmt.Println("Starting server at :8080")
	router := routes.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
