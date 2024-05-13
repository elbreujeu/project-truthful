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

	corsConfig := cors.DefaultConfig()
	// TODO: Not safe for production, change to specific origins
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	// sets router to use corsConfig
	router.Use(cors.New(corsConfig))

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
