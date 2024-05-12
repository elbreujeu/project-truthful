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
	token, err := token.GenerateJWT(int(id))
	if err != nil {
		log.Printf("Error while generating token: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error while generating token after registration",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created",
		"id":      id,
		"token":   token,
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

	// parses 2 query parameters, "count" and "start"
	countStr := c.Query("count")
	startStr := c.Query("start")

	count, err := basicfuncs.ConvertQueryParameterToInt(countStr, 10)
	if err != nil {
		log.Printf("Error while parsing query parameter: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error while parsing query parameter",
			"error":   err.Error(),
		})
		return
	}
	start, err := basicfuncs.ConvertQueryParameterToInt(startStr, 0)
	if err != nil {
		log.Printf("Error while parsing query parameter: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error while parsing query parameter",
			"error":   err.Error(),
		})
		return
	}

	if count < 0 || start < 0 {
		log.Printf("Error while parsing query parameter: negative values\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error while parsing query parameter",
			"error":   "negative values",
		})
		return
	}

	if count > 30 {
		count = 30
	}

	user, code, err := client.GetUserProfile(username, count, start)
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

func likeAnswer(c *gin.Context) {
	log.Printf("Received request to like answer from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.LikeAnswerInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	if infos.Like {
		code, err := client.LikeAnswer(requesterId, infos.AnswerId)
		if err != nil {
			log.Printf("Error while liking answer: %s\n", err.Error())
			c.JSON(code, gin.H{"message": "error while liking answer", "error": err.Error()})
			return
		}
	} else {
		code, err := client.UnlikeAnswer(requesterId, infos.AnswerId)
		if err != nil {
			log.Printf("Error while unliking answer: %s\n", err.Error())
			c.JSON(code, gin.H{"message": "error while unliking answer", "error": err.Error()})
			return
		}
	}

	if infos.Like {
		c.JSON(http.StatusCreated, gin.H{"message": "answer liked"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "answer unliked"})
	}
}

func deleteAnswer(c *gin.Context) {
	log.Printf("Received request to delete answer from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.DeleteAnswerInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	code, err := client.MarkAnswerAsDeleted(requesterId, infos.AnswerId)
	if err != nil {
		log.Printf("Error while deleting answer: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while deleting answer", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "answer deleted"})
}

func deleteQuestion(c *gin.Context) {
	log.Printf("Received request to delete question from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.DeleteQuestionInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	code, err := client.MarkQuestionAsDeleted(requesterId, infos.QuestionId)
	if err != nil {
		log.Printf("Error while deleting question: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while deleting question", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "question deleted"})
}

func updateUser(c *gin.Context) {
	// NOTE : In the future, change this function to a lot of smaller functions with PATCH requests
	log.Printf("Received request to update user from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.UpdateUserInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	code, err := client.UpdateUserInformations(requesterId, infos.DisplayName, infos.Email)
	if err != nil {
		log.Printf("Error while updating user: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while updating user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func promoteUser(c *gin.Context) {
	log.Printf("Received request to promote user from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.PromoteUserInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	code, err := client.PromoteUser(requesterId, infos.UserId, infos.PromoteType)
	if err != nil {
		log.Printf("Error while promoting user: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while promoting user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted"})
}

func moderationGetUserQuestions(c *gin.Context) {
	log.Printf("Received request to get user questions from ip %s\n", c.ClientIP())

	username := c.Param("user")

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

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	questions, code, err := client.ModerationGetUserQuestions(requesterId, username, start, count)
	if err != nil {
		log.Printf("Error while getting user questions: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while getting user questions", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func banUser(c *gin.Context) {
	log.Printf("Received request to ban user from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.BanUserInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	banId, code, err := client.BanUser(infos.UserId, requesterId, infos.Duration, infos.Reason)
	if err != nil {
		log.Printf("Error while banning user: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while banning user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user banned", "ban_id": banId})
}

func pardonUser(c *gin.Context) {
	log.Printf("Received request to pardon user from ip %s\n", c.ClientIP())

	requesterId, _, err := parseAndVerifyAccessToken(c)
	if err != nil {
		return
	}

	var infos models.PardonUserInfos
	err = c.ShouldBindJSON(&infos)
	if err != nil {
		log.Printf("Error while parsing request body: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while parsing request body", "error": err.Error()})
		return
	}

	pardonId, code, err := client.PardonUser(infos.BanId, requesterId)
	if err != nil {
		log.Printf("Error while pardoning user: %s\n", err.Error())
		c.JSON(code, gin.H{"message": "error while pardoning user", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user pardoned", "pardon_id": pardonId})
}

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
	r.POST("/like_answer", likeAnswer)
	r.POST("/delete_answer", deleteAnswer)
	r.POST("/delete_question", deleteQuestion)
	r.PUT("/users/update", updateUser)
	r.POST("/moderation/promote", promoteUser)
	r.GET("/moderation/get_user_questions/:user", moderationGetUserQuestions)
	r.POST("/moderation/ban_user", banUser)
	r.POST("/moderation/pardon_user", pardonUser)
}
