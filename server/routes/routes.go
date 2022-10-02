package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_truthful/client"
	"project_truthful/models"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "Hello world !"}`)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to create user from ip %s\n", r.RemoteAddr)
	// parses the request body and returns a models.CreateUserInfos struct
	// if the request body is not valid, returns an error
	var infos models.CreateUserInfos
	err := json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}
	// checks if all the fields are filled
	if infos.Username == "" || infos.Password == "" || infos.Email == "" || infos.Birthdate == "" {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	id, code, err := client.CreateUser(infos)
	if err != nil {
		log.Printf("Error while creating user: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while creating user", "error": "%s"}`, err.Error())
		return
	}
	log.Printf("User created with id %d\n", id)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "User created", "id": %d}`, id)
}

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/hello_world", homePage).Methods("GET")
	r.HandleFunc("/create_user", createUser).Methods("POST")
}
