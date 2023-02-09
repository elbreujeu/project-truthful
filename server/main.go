package main

import (
	"log"
	"os"
	"project_truthful/client/database"
	"project_truthful/client/token"
	"project_truthful/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const DEFAULT_PORT = "8080"

func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Printf("PORT not found in env, using default port %s\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}

	// Create a new Gin engine
	router := gin.Default()

	// Enable CORS and allow all origins
	// TODO: Research it and change it to a more secure way
	router.Use(cors.Default())

	// Add middleware
	routes.SetMiddleware(router)

	// Setup routes
	routes.SetupRoutes(router)

	var err error
	database.DB, err = database.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = token.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	log.Println("Starting server...")
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server shutted down.")
}
