package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project_truthful/client"
	"project_truthful/client/basicfuncs"
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
	fmt.Fprintf(w, `{"message": "Token refreshed", "token": "%s"}`, newToken)
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

	requesterId, _, err := parseAndVerifyAccessToken(w, r)
	if err != nil {
		return
	}

	if r.Body == nil {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	var infos models.FollowUserInfos
	err = json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	} else if infos.UserId == 0 {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	var message string
	var code int
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

func askQuestion(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to ask question from ip %s\n", r.RemoteAddr)

	accessToken, _, err := token.ParseAccessToken(r)
	requesterId := 0
	var code int
	if err == nil {
		requesterId, code, err = token.VerifyJWT(accessToken)
		if err != nil {
			log.Printf("Error while checking token: %s\n", err.Error())
			w.WriteHeader(code)
			fmt.Fprintf(w, `{"message": "error while checking token", "error": "%s"}`, err.Error())
			return
		}
	}

	if r.Body == nil {
		log.Printf("Error while parsing request body: missing fields\n")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "missing fields"}`)
		return
	}

	var infos models.AskQuestionInfos
	err = json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	id, code, err := client.AskQuestion(infos.QuestionText, requesterId, r.RemoteAddr, infos.IsAuthorAnonymous, infos.UserId)
	if err != nil {
		log.Printf("Error while asking question: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while asking question", "error": "%s"}`, err.Error())
		return
	}

	log.Printf("Question asked with id %d\n", id)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "Question asked", "id": %d}`, id)
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to get questions from ip %s\n", r.RemoteAddr)

	//get the "count" parameter from query parameters
	countStr := r.URL.Query().Get("count")
	startStr := r.URL.Query().Get("start")

	count, err := basicfuncs.ConvertQueryParameterToInt(countStr, 10)
	if err != nil {
		log.Printf("Error while parsing count: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid count", "error": "%s"}`, err.Error())
		return
	}

	start, err := basicfuncs.ConvertQueryParameterToInt(startStr, 0)
	if err != nil {
		log.Printf("Error while parsing count: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid count", "error": "%s"}`, err.Error())
		return
	}

	userId, _, err := parseAndVerifyAccessToken(w, r)
	if err != nil {
		return
	}
	questions, code, err := client.GetQuestions(userId, start, count)
	if err != nil {
		log.Printf("Error while getting questions: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while getting questions", "error": "%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(questions)
	if err != nil {
		log.Printf("Error while encoding questions: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message": "error while encoding questions", "error": "%s"}`, err.Error())
		return
	}
}

func answerQuestion(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to answer question from ip %s\n", r.RemoteAddr)

	requesterId, _, err := parseAndVerifyAccessToken(w, r)
	if err != nil {
		return
	}

	var infos models.AnswerQuestionInfos
	err = json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	id, code, err := client.AnswerQuestion(requesterId, infos.QuestionId, infos.AnswerText, r.RemoteAddr)
	if err != nil {
		log.Printf("Error while answering question: %s\n", err.Error())
		w.WriteHeader(code)
		fmt.Fprintf(w, `{"message": "error while answering question", "error": "%s"}`, err.Error())
		return
	}

	log.Printf("Question answered with id %d\n", id)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "Question answered", "id": %d}`, id)
}

func likeAnswer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to like answer from ip %s\n", r.RemoteAddr)

	requesterId, _, err := parseAndVerifyAccessToken(w, r)
	if err != nil {
		return
	}

	var infos models.LikeAnswerInfos
	err = json.NewDecoder(r.Body).Decode(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid request body", "error": "%s"}`, err.Error())
		return
	}

	if infos.Like {
		code, err := client.LikeAnswer(requesterId, infos.AnswerId)
		if err != nil {
			log.Printf("Error while liking answer: %s\n", err.Error())
			w.WriteHeader(code)
			fmt.Fprintf(w, `{"message": "error while liking answer", "error": "%s"}`, err.Error())
			return
		}
	} else {
		code, err := client.UnlikeAnswer(requesterId, infos.AnswerId)
		if err != nil {
			log.Printf("Error while unliking answer: %s\n", err.Error())
			w.WriteHeader(code)
			fmt.Fprintf(w, `{"message": "error while unliking answer", "error": "%s"}`, err.Error())
			return
		}
	}

	log.Printf("Answer liked by user %d\n", requesterId)
	w.WriteHeader(http.StatusOK)
	if infos.Like {
		fmt.Fprintf(w, `{"message": "Answer liked"}`)
	} else {
		fmt.Fprintf(w, `{"message": "Answer unliked"}`)
	}
}

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/hello_world", homePage).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/refresh_token", refreshToken).Methods("GET")
	r.HandleFunc("/get_user_profile/{user}", getUserProfile).Methods("GET")
	r.HandleFunc("/follow_user", followUser).Methods("POST")
	r.HandleFunc("/ask_question", askQuestion).Methods("POST")
	r.HandleFunc("/get_questions", getQuestions).Methods("GET")
	r.HandleFunc("/answer_question", answerQuestion).Methods("POST")
	r.HandleFunc("/like_answer", likeAnswer).Methods("POST")
}
