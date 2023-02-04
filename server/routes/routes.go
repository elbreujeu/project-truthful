package routes

import (
	"log"
	"net/http"
	"project_truthful/client"
	"project_truthful/client/basicfuncs"
	"project_truthful/client/token"
	"project_truthful/models"

	"github.com/gin-gonic/gin"
)

func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello world !"})
}

func register(c *gin.Context) {
	log.Printf("Received request to create user from ip %s\n", c.ClientIP())

	var infos models.RegisterInfos
	if err := c.ShouldBindJSON(&infos); err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if infos.Username == "" || infos.Password == "" || infos.Email == "" || infos.Birthdate == "" {
		log.Printf("Error while parsing request body: missing fields\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   "missing fields",
		})
		return
	}

	id, code, err := client.Register(infos)
	if err != nil {
		log.Printf("Error while creating user: %s\n", err.Error())
		c.JSON(code, gin.H{
			"message": "error while creating user",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("User created with id %d\n", id)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"id":      id,
	})
}

func login(c *gin.Context) {
	log.Printf("Received request to login from ip %s\n", c.ClientIP())

	var infos models.LoginInfos
	err := c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if infos.Username == "" || infos.Password == "" {
		log.Printf("Error while parsing request body: missing fields\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   "missing fields",
		})
		return
	}

	token, code, err := client.Login(infos)
	if err != nil {
		log.Printf("Error while logging in: %s\n", err.Error())
		c.JSON(code, gin.H{
			"message": "error while logging in",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("User logged in with token %s\n", token)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in",
		"token":   token,
	})
}

func refreshToken(c *gin.Context) {
	log.Printf("Received request to refresh token from ip %s\n", c.ClientIP())

	accessToken, code, err := token.ParseAccessToken(c)
	if err != nil {
		log.Printf("Error while parsing token: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while parsing token", "error": err.Error()})
		return
	}

	newToken, code, err := token.RefreshJWT(accessToken)
	if err != nil {
		log.Printf("Error while checking token: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while checking token", "error": err.Error()})
		return
	}
	log.Printf("Token refreshed")
	c.JSON(code, gin.H{"message": "Token refreshed", "token": newToken})
}

func getUserProfile(c *gin.Context) {
	log.Printf("Received request to get user from ip %s\n", c.ClientIP())

	username := c.Param("user")

	user, code, err := client.GetUserProfile(username)
	if err != nil {
		log.Printf("Error while getting user: %s\n", err.Error())
		c.JSON(code, gin.H{
			"message": "error while getting user",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("User %s found\n", username)
	c.JSON(http.StatusOK, user)
}

func followUser(c *gin.Context) {
	log.Printf("Received request to follow user from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.FollowUserInfos
	if err := c.ShouldBindJSON(&infos); err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	} else if infos.UserId == 0 {
		log.Printf("Error while parsing request body: missing fields\n")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": "missing fields"})
		return
	}

	var message string
	if infos.Follow {
		_, err = client.FollowUser(requesterId, infos.UserId)
		message = "User followed"
	} else {
		_, err = client.UnfollowUser(requesterId, infos.UserId)
		message = "User unfollowed"
	}
	if err != nil {
		log.Printf("Error while following user: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error while following user", "error": err.Error()})
		return
	}

	log.Printf("User %d followed user %d\n", requesterId, infos.UserId)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func askQuestion(c *gin.Context) {
	log.Printf("Received request to ask question from ip %s\n", c.ClientIP())
	accessToken, _, err := token.ParseAccessToken(c)
	requesterId := 0
	var code int
	if err == nil {
		requesterId, code, err = token.VerifyJWT(accessToken)
		if err != nil {
			log.Printf("Error while checking token: %s\n", err.Error())
			c.JSON(code, gin.H{"message": "error while checking token", "error": err.Error()})
			return
		}
	}

	var infos models.AskQuestionInfos
	if err = c.ShouldBindJSON(&infos); err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body", "error": err.Error()})
		return
	}

	id, code, err := client.AskQuestion(infos.QuestionText, requesterId, c.ClientIP(), infos.IsAuthorAnonymous, infos.UserId)
	if err != nil {
		log.Printf("Error while asking question: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while asking question", "error": err.Error()})
		return
	}

	log.Printf("Question asked with id %d\n", id)
	c.JSON(http.StatusCreated, gin.H{"message": "Question asked", "id": id})
}

func getQuestions(c *gin.Context) {
	log.Printf("Received request to get questions from ip %s\n", c.ClientIP())

	//get the "count" parameter from query parameters
	countStr := c.Query("count")
	startStr := c.Query("start")

	count, err := basicfuncs.ConvertQueryParameterToInt(countStr, 10)
	if err != nil {
		log.Printf("Error while parsing count: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid count",
			"error":   err.Error(),
		})
		return
	}

	start, err := basicfuncs.ConvertQueryParameterToInt(startStr, 0)
	if err != nil {
		log.Printf("Error while parsing count: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid count",
			"error":   err.Error(),
		})
		return
	}

	userId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	questions, code, err := client.GetQuestions(userId, start, count)
	if err != nil {
		log.Printf("Error while getting questions: %s\n", err.Error())
		c.JSON(code, gin.H{
			"message": "error while getting questions",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func answerQuestion(c *gin.Context) {
	log.Printf("Received request to answer question from ip %s\n", c.ClientIP())
	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.AnswerQuestionInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	id, code, err := client.AnswerQuestion(requesterId, infos.QuestionId, infos.AnswerText, c.ClientIP())
	if err != nil {
		log.Printf("Error while answering question: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while answering question", "error": err.Error()})
		return
	}

	log.Printf("Question answered with id %d\n", id)
	c.JSON(http.StatusCreated, gin.H{"message": "question answered", "id": id})
}

// func likeAnswer(w http.ResponseWriter, r *http.Request) {
// 	log.Printf("Received request to like answer from ip %s\n", r.RemoteAddr)

// 	requesterId, _, err := parseAndVerifyAccessToken(w, r)
// 	if err != nil {
// 		return
// 	}

// 	var infos models.LikeAnswerInfos
// 	err = json.NewDecoder(r.Body).Decode(&infos)
// 	if err != nil {
// 		log.Printf("Error while parsing request body: %s\n", err.Error())
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintf(w, `{"message": "invalid request body", "error": "%s"}`, err.Error())
// 		return
// 	}

// 	if infos.Like {
// 		code, err := client.LikeAnswer(requesterId, infos.AnswerId)
// 		if err != nil {
// 			log.Printf("Error while liking answer: %s\n", err.Error())
// 			w.WriteHeader(code)
// 			fmt.Fprintf(w, `{"message": "error while liking answer", "error": "%s"}`, err.Error())
// 			return
// 		}
// 	} else {
// 		code, err := client.UnlikeAnswer(requesterId, infos.AnswerId)
// 		if err != nil {
// 			log.Printf("Error while unliking answer: %s\n", err.Error())
// 			w.WriteHeader(code)
// 			fmt.Fprintf(w, `{"message": "error while unliking answer", "error": "%s"}`, err.Error())
// 			return
// 		}
// 	}

// 	log.Printf("Answer liked by user %d\n", requesterId)
// 	w.WriteHeader(http.StatusOK)
// 	if infos.Like {
// 		fmt.Fprintf(w, `{"message": "Answer liked"}`)
// 	} else {
// 		fmt.Fprintf(w, `{"message": "Answer unliked"}`)
// 	}
// }

func SetupRoutes(r *gin.Engine) {
	r.GET("/hello_world", helloWorld)
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/refresh_token", refreshToken)
	r.GET("/get_user_profile/:user", getUserProfile)
	r.POST("/follow_user", followUser)
	r.POST("/ask_question", askQuestion)
	r.GET("/get_questions", getQuestions)
	r.POST("/answer_question", answerQuestion)
	// r.POST("/like_answer", likeAnswer)
}
