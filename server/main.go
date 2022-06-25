package main

import (
	"log"
	"net/http"
	"os"
	"project_truthful/client/database"
	"project_truthful/routes"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const DEFAULT_PORT = "8080"

func main() {
	err := godotenv.Overload("../.env")
	if err != nil {
		log.Fatal("Error overloading .env file")
	}
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Printf("PORT not found in env, using default port %s\n", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	router := mux.NewRouter()
	routes.SetMiddleware(router)
	routes.SetupRoutes(router)

	err = database.Init()
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(router),
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	waiter := make(chan error)
	log.Println("Starting server ...")
	go func() {
		err := s.ListenAndServe()
		waiter <- err
	}()
	log.Printf("Project Truthful server running on http://localhost:%s\n", port)
	err = <-waiter
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server shutted down.")
}
