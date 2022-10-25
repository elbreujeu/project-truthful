package routes

import (
	"fmt"
	"log"
	"net/http"
	"project_truthful/client/token"

	"github.com/gorilla/mux"
)

func SetMiddleware(r *mux.Router) {
	r.Use(setCORS)
	r.Use(setJSONResponse)
}

func setCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func setJSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func parseAndVerifyAccessToken(w http.ResponseWriter, r *http.Request) (int, int, error) {
	accessToken, code, err := token.ParseAccessToken(r)
	if err != nil {
		log.Printf("Error while parsing token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while parsing token", "error": "%s"}`, err.Error())
		return 0, code, err
	}

	requesterId, code, err := token.VerifyJWT(accessToken)
	if err != nil {
		log.Printf("Error while checking token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while checking token", "error": "%s"}`, err.Error())
		return 0, code, err
	}
	return requesterId, http.StatusOK, nil
}
