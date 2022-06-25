package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "Hello world !"}`)
}

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/hello_world", homePage).Methods("GET")
}
