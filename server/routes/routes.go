package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_truthful/client"
	"project_truthful/client/token"
	"project_truthful/models"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message": "Hello world !"}`)
}

func register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to create user from ip %s\n", r.RemoteAddr)

	var infos models.RegisterInfos
	err := json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	if infos.Username == "" || infos.Password == "" || infos.Email == "" || infos.Birthdate == "" {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	id, code, err := client.Register(infos)
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

func login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to login from ip %s\n", r.RemoteAddr)

	var infos models.LoginInfos
	err := json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	if infos.Username == "" || infos.Password == "" {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	token, code, err := client.Login(infos)
	if err != nil {
		log.Printf("Error while logging in: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while logging in", "error": "%s"}`, err.Error())
		return
	}
	log.Printf("User logged in with token %s\n", token)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "User logged in", "token": "%s"}`, token)
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to refreshed token from ip %s\n", r.RemoteAddr)

	accessToken, code, err := token.ParseAccessToken(r)
	if err != nil {
		log.Printf("Error while parsing token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while parsing token", "error": "%s"}`, err.Error())
		return
	}

	newToken, code, err := token.RefreshJWT(accessToken)
	if err != nil {
		log.Printf("Error while checking token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while checking token", "error": "%s"}`, err.Error())
		return
	}
	log.Printf("Token refresheded")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"message": "Token refresheded", "token": "%s"}`, newToken)
}

func getUserProfile(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to get user from ip %s\n", r.RemoteAddr)

	vars := mux.Vars(r)
	username := vars["user"]

	user, code, err := client.GetUserProfile(username)
	if err != nil {
		log.Printf("Error while getting user: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while getting user", "error": "%s"}`, err.Error())
		return
	}
	log.Printf("User %s found\n", username)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("Error while encoding user: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message": "error while encoding user", "error": "%s"}`, err.Error())
		return
	}
}

func followUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to follow user from ip %s\n", r.RemoteAddr)

	accessToken, code, err := token.ParseAccessToken(r)
	if err != nil {
		log.Printf("Error while parsing token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while parsing token", "error": "%s"}`, err.Error())
		return
	}

	requesterId, code, err := token.VerifyJWT(accessToken)
	if err != nil {
		log.Printf("Error while checking token: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while checking token", "error": "%s"}`, err.Error())
		return
	}

	var infos models.FollowUserInfos
	err = json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	var message string
	if infos.Follow {
		code, err = client.FollowUser(requesterId, infos.UserId)
		message = "User followed"
	} else {
		code, err = client.UnfollowUser(requesterId, infos.UserId)
		message = "User unfollowed"
	}
	if err != nil {
		log.Printf("Error while following user: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while following user", "error": "%s"}`, err.Error())
		return
	}

	log.Printf("User %d followed user %d\n", requesterId, infos.UserId)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "%s"}`, message)
}

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/hello_world", homePage).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/refresh_token", refreshToken).Methods("GET")
	r.HandleFunc("/{user}", getUserProfile).Methods("GET")
	r.HandleFunc("/follow_user", followUser).Methods("POST")
}
